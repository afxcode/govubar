import Clutter from 'gi://Clutter';
import GObject from 'gi://GObject';
import St from 'gi://St';
import GLib from 'gi://GLib';
import Gio from 'gi://Gio';

import {Extension, gettext as _} from 'resource:///org/gnome/shell/extensions/extension.js';
import * as PanelMenu from 'resource:///org/gnome/shell/ui/panelMenu.js';
import * as Main from 'resource:///org/gnome/shell/ui/main.js';

const Indicator = GObject.registerClass(
    class Indicator extends PanelMenu.Button {
        _init() {
            super._init(0.0, _('GoVUBar'));

            this._destroyed = false;
            this._subprocess = null;

            this.label = new St.Label({
                style_class: 'govubar-panel-label',
                text: "○○○○○",
                y_expand: true,
                y_align: Clutter.ActorAlign.CENTER
            });

            this.add_child(this.label);

            this._startSubprocess();
        }

        _startSubprocess() {
            try {
                this._subprocess = Gio.Subprocess.new(
                    ['bash', '-c', 'govubar'],
                    Gio.SubprocessFlags.STDOUT_PIPE | Gio.SubprocessFlags.HAVE_STDERR
                );

                this._stdout = new Gio.DataInputStream({
                    base_stream: this._subprocess.get_stdout_pipe(),
                    close_base_stream: true,
                });

                this._readOutput();
            } catch (e) {
                logError(e, 'Failed to start govubar');
            }
        }

        _readOutput() {
            if (this._destroyed || !this._stdout)
                return;

            this._stdout.read_line_async(GLib.PRIORITY_DEFAULT, null, (stream, result) => {
                    try {
                        if (this._destroyed) return;

                        const [line] = stream.read_line_finish_utf8(result);

                        if (line === null) return;

                        if (this.label) this.label.text = line;

                        this._readOutput();
                    } catch (e) {
                        if (!this._destroyed) logError(e);
                    }
                }
            );
        }

        destroy() {
            this._destroyed = true;

            if (this._subprocess) {
                try {
                    this._subprocess.force_exit();
                } catch (e) {
                    logError(e);
                }

                this._subprocess = null;
            }

            this._stdout = null;
            this.label = null;

            super.destroy();
        }
    });

export default class IndicatorExampleExtension extends Extension {
    enable() {
        this._indicator = new Indicator();
        Main.panel.addToStatusArea(this.uuid, this._indicator);
    }

    disable() {
        if (this._indicator) {
            this._indicator.destroy();
            this._indicator = null;
        }
    }
}

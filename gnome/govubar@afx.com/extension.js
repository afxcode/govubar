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
        const script = `govubar`;

        try {
            this._subprocess = Gio.Subprocess.new(
                ['bash', '-c', script],
                Gio.SubprocessFlags.STDOUT_PIPE | Gio.SubprocessFlags.HAVE_STDERR
            );

            const stdoutStream = new Gio.DataInputStream({
                base_stream: this._subprocess.get_stdout_pipe(),
                close_base_stream: true,
            });

            this._readOutput(stdoutStream);
        } catch (e) {
            logError(e);
        }
    }

    _readOutput(stdout) {
        if (!this.label || this.label.isDestroyed()) return;

        stdout.read_line_async(GLib.PRIORITY_LOW, null, (stream, result) => {
            try {
                const [line] = stream.read_line_finish_utf8(result);

                if (line !== null) {
                    this.label.text = line;
                    this._readOutput(stdout);
                }
            } catch (e) {
                logError(e);
            }
        });
    }

    vfunc_dispose() {
        super.vfunc_dispose();

        if (this._subprocess) {
            this._subprocess.force_exit();
            this._subprocess = null;
        }

        if (this.label) {
            this.label.destroy();
            this.label = null;
        }
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


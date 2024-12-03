package main

import (
	"context"
	"strings"
	"sync/atomic"

	"github.com/noriah/catnip/dsp"
	"github.com/noriah/catnip/util"
)

// Constants
const (
	// ScalingWindow in seconds
	ScalingWindow = 1.5
	// PeakThreshold is the threshold to not draw if the peak is less.
	PeakThreshold = 0.001
)

// Writer handles drawing our visualizer.
type Writer struct {
	Smoother  dsp.Smoother
	trackZero int
	binCount  int
	window    *util.MovingWindow
}

func NewWriter() *Writer {
	return &Writer{
		binCount: 1,
	}
}

// Init initializes the display.
// Should be called before any other display method.
func (d *Writer) Init(sampleRate float64, sampleSize int) error {
	windowSize := ((int(ScalingWindow * sampleRate)) / sampleSize) * 2
	d.window = util.NewMovingWindow(windowSize)
	return nil
}

// Close will stop display and clean up the terminal.
func (d *Writer) Close() error {
	return nil
}

func (d *Writer) SetBinCount(int) {}

func (d *Writer) SetInvertDraw(bool) {
}

// Start display is bad.
func (d *Writer) Start(ctx context.Context) context.Context {
	return ctx
}

// Stop display not work.
func (d *Writer) Stop() error {
	return nil
}

var setBarGraphIsRun atomic.Bool

// Draw takes data and draws.
func (d *Writer) Write(buffers [][]float64, channels int) error {
	peak := 0.0
	bins := d.Bins(channels)

	for i := 0; i < channels; i++ {
		for _, val := range buffers[i][:bins] {
			if val > peak {
				peak = val
			}
		}
	}

	scale := 1.0

	if peak >= PeakThreshold {
		d.trackZero = 0

		// do some scaling if we are above the PeakThreshold
		d.window.Update(peak)

	} else {
		if d.trackZero++; d.trackZero == 5 {
			d.window.Recalculate()
		}
	}

	vMean, vSD := d.window.Stats()

	if t := vMean + (2.0 * vSD); t > 1.0 {
		scale = t
	}

	scale = 100.0 / scale

	ch := 0
	total := 0.0
	for xSet, chBins := range buffers {
		for xBar := 0; xBar < d.binCount; xBar++ {
			xBin := (xBar * (1 - xSet)) + (((d.binCount - 1) - xBar) * xSet)

			total += chBins[xBin] * scale
			ch++
		}
	}

	go func() {
		if setBarGraphIsRun.Load() {
			return
		}
		setBarGraphIsRun.Store(true)
		defer setBarGraphIsRun.Store(false)

		setDisplayBarGraph(int(total / float64(ch)))
	}()
	return nil
}

// Bins returns the number of bars we will draw.
func (d *Writer) Bins(chCount int) int {
	return d.binCount
}

func generateBarGraph(value int) string {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	totalBars := 5 // Total segments in the bar
	filledBars := value * totalBars / 100

	return strings.Repeat("○", totalBars-filledBars) + strings.Repeat("●", filledBars)
}

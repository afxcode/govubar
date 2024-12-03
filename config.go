package main

import (
	"errors"

	"github.com/noriah/catnip/dsp"
	"github.com/noriah/catnip/input"
)

// Config is a temporary struct to define parameters
type config struct {
	// Backend is the backend name from list-backends
	backend string
	// Device is the device name from list-devices
	device string
	// SampleRate is the rate at which samples are read
	sampleRate float64
	// SmoothFactor factor of smooth
	smoothFactor float64
	// Smoothing method used to do time smoothing.
	smoothingMethod int
	// SampleSize is how much we draw. Play with it
	sampleSize int
}

// NewZeroConfig returns a zero config
// it is the "default"
func newZeroConfig() config {
	return config{
		backend:         input.DefaultBackend(),
		sampleRate:      44100,
		sampleSize:      512,
		smoothFactor:    1,
		smoothingMethod: int(dsp.SmoothAverage),
	}
}

// Sanitize cleans things up
func (cfg *config) validate() error {
	if cfg.sampleRate < float64(cfg.sampleSize) {
		return errors.New("sample rate lower than sample size")
	}

	if cfg.sampleSize < 4 {
		return errors.New("sample size too small (4+ required)")
	}

	switch {
	case cfg.smoothFactor > 99.99:
		cfg.smoothFactor = 0.9999
	case cfg.smoothFactor < 0.00001:
		cfg.smoothFactor = 0.00001
	default:
		cfg.smoothFactor /= 100.0
	}
	return nil
}

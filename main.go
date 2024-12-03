package main

import (
	"context"
	"fmt"
	"log"

	"github.com/integrii/flaggy"
	"github.com/noriah/catnip"
	"github.com/noriah/catnip/dsp"
	"github.com/noriah/catnip/dsp/window"
	"github.com/noriah/catnip/input"
	_ "github.com/noriah/catnip/input/all"
)

const (
	AppName = "GoVUBar"
	AppDesc = "Audio VU Meter on the Top Bar"
	Version = "0.1.0"
)

func main() {
	cfg := newZeroConfig()

	if doFlags(&cfg) {
		return
	}

	if err := cfg.validate(); err != nil {
		log.Fatalln("invalid config: ", err)
	}

	smoother := dsp.NewSmoother(dsp.SmootherConfig{
		SampleSize:      cfg.sampleSize,
		SampleRate:      cfg.sampleRate,
		ChannelCount:    2,
		SmoothingFactor: cfg.smoothFactor,
		SmoothingMethod: dsp.SmoothingMethod(cfg.smoothingMethod),
	})

	writer := NewWriter()
	writer.Init(cfg.sampleRate, cfg.sampleSize)

	catnipCfg := catnip.Config{
		Backend:      cfg.backend,
		Device:       cfg.device,
		SampleRate:   cfg.sampleRate,
		SampleSize:   cfg.sampleSize,
		ChannelCount: 2,
		UseThreaded:  false,
		SetupFunc:    func() error { return nil },
		StartFunc: func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		},
		CleanupFunc: func() error { return nil },
		Output:      writer,
		Windower:    window.Lanczos(),
		Analyzer: dsp.NewAnalyzer(dsp.AnalyzerConfig{
			SampleRate:    cfg.sampleRate,
			SampleSize:    cfg.sampleSize,
			SquashLow:     true,
			SquashLowOld:  true,
			DontNormalize: false,
			BinMethod:     dsp.MaxSampleValue(),
		}),
		Smoother: smoother,
	}

	wg, ctx := getDisplayContext()
	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := catnip.Run(&catnipCfg, ctx); err != nil {
			log.Println("Capture audio failed: ", err)
		}
	}()

	runDisplay()
}

func doFlags(cfg *config) bool {
	parser := flaggy.NewParser(AppName)
	parser.Description = AppDesc
	parser.Version = Version

	listBackendsCmd := flaggy.Subcommand{
		Name:                 "list-backends",
		ShortName:            "lb",
		Description:          "list all supported backends",
		AdditionalHelpAppend: "\nuse the full name after the '-'",
	}

	parser.AttachSubcommand(&listBackendsCmd, 1)

	listDevicesCmd := flaggy.Subcommand{
		Name:                 "list-devices",
		ShortName:            "ld",
		Description:          "list all devices for a backend",
		AdditionalHelpAppend: "\nuse the full name after the '-'",
	}

	parser.AttachSubcommand(&listDevicesCmd, 1)

	parser.String(&cfg.backend, "b", "backend", "backend name")
	parser.String(&cfg.device, "d", "device", "device name")
	parser.Float64(&cfg.sampleRate, "r", "rate", "sample rate")
	parser.Int(&cfg.sampleSize, "n", "samples", "sample size")
	parser.Float64(&cfg.smoothFactor, "sf", "smoothing", "smooth factor (0-100)")
	parser.Int(&cfg.smoothingMethod, "sm", "smooth-method", "smoothing method (0, 1, 2, 3, 4, 5)")

	if err := parser.Parse(); err != nil {
		log.Fatalln("failed to parse arguments:", err)
	}

	switch {
	case listBackendsCmd.Used:
		defaultBackend := input.DefaultBackend()

		fmt.Println("all backends. '*' marks default")

		for _, backend := range input.Backends {
			star := ' '
			if defaultBackend == backend.Name {
				star = '*'
			}

			fmt.Printf("- %s %c\n", backend.Name, star)
		}
		return true

	case listDevicesCmd.Used:
		backend, err := input.InitBackend(cfg.backend)
		if err != nil {
			log.Fatalln("failed to initialize backend:", err)
		}

		devices, err := backend.Devices()
		if err != nil {
			log.Fatalln("failed to get devices:", err)
		}

		// We don't really need the default device to be indicated.
		defaultDevice, _ := backend.DefaultDevice()

		fmt.Printf("all devices for %q backend. '*' marks default\n", cfg.backend)

		for idx := range devices {
			star := ' '
			if defaultDevice != nil && devices[idx].String() == defaultDevice.String() {
				star = '*'
			}

			fmt.Printf("- %v %c\n", devices[idx], star)
		}
		return true
	}
	return false
}

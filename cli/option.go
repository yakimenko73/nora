package cli

import (
	"errors"
	"github.com/illatior/nora/lib/metric"
	"github.com/illatior/nora/cli/screen"
	"github.com/mum4k/termdash/container"
	"time"
)

type Option interface {
	apply(ui *cui) error
}

type option func(ui *cui) error

func (o option) apply(ui *cui) error {
	return o(ui)
}

func WithChangeDisplayableIntervalDelta(d time.Duration) Option {
	return option(func(ui *cui) error {
		if d <= 0 {
			return errors.New("`change displayable interval delta` should be > 0")
		}

		ui.changeDisplayableIntervalDelta = d
		return nil
	})
}

func WithUpdateInterval(d time.Duration) Option {
	return option(func(ui *cui) error {
		if d <= 0 {
			return errors.New("update interval should be > 0")
		}

		ui.updateInterval = d
		return nil
	})
}

func WithoutDefaultScreens() Option {
	return option(func(ui *cui) error {
		ui.screens = nil
		return nil
	})
}

func WithDebugScreen() Option {
	return option(func(ui *cui) error {
		debugScreen, err := screen.NewDebugScreen()
		if err != nil {
			return err
		}

		ui.screens = append(ui.screens, debugScreen)
		return nil
	})
}

func WithCustomScreen(s screen.Screen) Option {
	return option(func(ui *cui) error {
		ui.screens = append(ui.screens, s)
		return nil
	})
}

func WithCustomContainer(c *container.Container) Option {
	return option(func(ui *cui) error {
		ui.c = c
		return nil
	})
}

func WithCustomMetrics(m metric.Metrics) Option {
	return option(func(ui *cui) error {
		ui.metrics = m
		return nil
	})
}

func WithCustomSubs(s subsFunc) Option {
	return option(func(ui *cui) error {
		ui.subs = s
		return nil
	})
}

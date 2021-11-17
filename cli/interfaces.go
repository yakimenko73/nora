package cli

import (
	"context"
	"github.com/illatior/nora/lib/metric"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/linechart"
)

type ConsoleUserInterface interface {
	Run(ctx context.Context, metrics <-chan *metric.Result, dispatchDone <-chan bool) error
	GetDoneChan() <-chan bool

	ChangeFullscreenState() error
	NextScreen() error
	PreviousScreen() error

	IsFullscreen() bool
}

type InfoPanel interface {
	widgetapi.Widget
	// FIXME
}

type Chart interface {
	widgetapi.Widget
	Series(label string, values []float64, opts ...linechart.SeriesOption) error
}

type FullscreenChart interface {
	Series(label string, values []float64, opts ...linechart.SeriesOption) error
}

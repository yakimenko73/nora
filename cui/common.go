package cui

import (
	"context"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/linechart"
)

const (
	SCREEN_ID = "scr"
)

type ConsoleUserInterface interface {
	Run(ctx context.Context, done chan<- bool) error
	AcceptMetric(m *metric.Result)

	ChangeFullscreenState() error
	NextScreen() error
	PreviousScreen() error

	IsFullscreen() bool
}

type Screen interface {
	GetBody() grid.Element

	GetHeader() grid.Element
	GetFooter() grid.Element
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

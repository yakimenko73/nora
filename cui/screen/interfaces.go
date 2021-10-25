package screen

import (
	"context"
	"github.com/illatior/nora/core/metric"
	"github.com/mum4k/termdash/container/grid"
	"time"
)

type Screen interface {
	GetBody() grid.Element
	GetHeader() grid.Element
	GetFooter() grid.Element

	ChangeDisplayInterval(t time.Duration)

	GetMetricsChan() chan<- metric.Metrics
	Run(ctx context.Context)
}

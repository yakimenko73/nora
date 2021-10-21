package screen

import (
	"context"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/mum4k/termdash/container/grid"
)

type Screen interface {
	GetBody() grid.Element
	GetHeader() grid.Element
	GetFooter() grid.Element

	GetMetricsChan() chan<- metric.Metrics
	Run(ctx context.Context)
}

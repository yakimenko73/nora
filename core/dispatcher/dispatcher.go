package dispatcher

import (
	"context"
	"load-testing/core/metric"
	"load-testing/core/worker"
)

type Dispatcher interface {
	AddWorker(id string, worker *worker.Worker) error
	Dispatch(ctx context.Context, metricConsumer *metric.MetricConsumer) error
}

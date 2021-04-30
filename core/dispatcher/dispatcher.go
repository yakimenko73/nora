package dispatcher

import (
	"context"
	"load-testing/core/job"
	"load-testing/core/metric"
)

type Dispatcher interface {
	AddJob(id string, job job.Job) error
	Dispatch(ctx context.Context, metricConsumer *metric.MetricConsumer) error
}

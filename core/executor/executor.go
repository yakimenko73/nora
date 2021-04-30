package executor

import (
	"context"
	"load-testing/core/job"
	"load-testing/core/metric"
)

type Executor interface {
	Execute(job job.Job, ctx context.Context, consumer metric.MetricConsumer) error
}


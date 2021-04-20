package worker

import (
	"context"
	"load-testing/core/job"
	"load-testing/core/metric"
)

type Worker interface {
	Run(ctx context.Context, metricConsumer *metric.MetricConsumer) error
	Stop()
	AddJob(job job.Job) error
}

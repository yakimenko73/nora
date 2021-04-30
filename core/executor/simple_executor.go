package executor

import (
	"context"
	"load-testing/config"
	"load-testing/core/job"
	"load-testing/core/metric"
)

type simpleExecutor struct {

}


func New(cfg *config.LoadTestConfig) Executor {
	return &simpleExecutor{

	}
}


func (e *simpleExecutor) Execute(job job.Job, ctx context.Context, consumer metric.MetricConsumer) error {
	return job.Complete()
}

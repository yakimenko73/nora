package executor

import (
	"context"
	"github.com/illatior/load-testing/core/job"
	"github.com/illatior/load-testing/core/metric"
	"sync"
)

type Executor interface {
	AddJob(job *job.Job)
	ScheduleExecution(ctx context.Context, wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result)
}

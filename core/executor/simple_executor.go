package executor

import (
	"context"
	"github.com/illatior/load-testing/core/job"
	"github.com/illatior/load-testing/core/metric"
	"sync"
	"time"
)

type simpleExecutor struct {
	jobs []*job.Job
}

func New() Executor {
	return &simpleExecutor{
		jobs: make([]*job.Job, 0),
	}
}

func (e *simpleExecutor) AddJob(j *job.Job) {
	e.jobs = append(e.jobs, j)
}

func (e *simpleExecutor) ScheduleExecution(ctx context.Context, wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			tickTime := time.Now()
			for _, j := range e.jobs {
				results <- (*j).Complete(tickTime)
			}
		}
	}
}

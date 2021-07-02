package executor

import (
	"context"
	"errors"
	"github.com/illatior/load-testing/core/job"
	"github.com/illatior/load-testing/core/metric"
	"sync"
	"sync/atomic"
	"time"
)

type roundRobinExecutor struct {
	jobs []*job.Job

	next int64
}

func NewRoundRobinExecutor() Executor {
	return &roundRobinExecutor{
		jobs: make([]*job.Job, 0),
		next: 0,
	}
}

func (e *roundRobinExecutor) AddJob(j *job.Job) {
	e.jobs = append(e.jobs, j)
}

func (e *roundRobinExecutor) ScheduleExecution(ctx context.Context, wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			next, err := e.getNext()
			if err != nil {
				panic(err)
			}

			results <- (*e.jobs[next]).Complete(time.Now())
		}
	}
}

func (e *roundRobinExecutor) getNext() (int64, error) {
	if len(e.jobs) == 0 {
		return 0, errors.New("Executor not initialized (jobs == 0)!")
	}

	atomic.AddInt64(&e.next, 1)
	if e.next == int64(len(e.jobs)) {
		e.next = 0
	}

	return e.next, nil
}
package dispatcher

import (
	"context"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/metric"
	"sync"
	"time"
)

type Scheduler interface {
	GetNextExecution(elapsed time.Duration, hits uint64) (next time.Duration, stop bool)
}

func Dispatch(ctx context.Context, scheduler Scheduler, executor executor.Executor, duration time.Duration, workers uint64) <-chan *metric.Result {
	var wg sync.WaitGroup

	ticks := make(chan interface{})
	results := make(chan *metric.Result)

	for i := uint64(0); i < workers; i++ {
		wg.Add(1)
		go executor.ScheduleExecution(&wg, ticks, results)
	}

	go func() {
		defer close(results)
		defer wg.Wait()
		defer close(ticks)

		start, executed := time.Now(), uint64(0)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}

			timeElapsed := time.Since(start)
			if timeElapsed > duration && duration != 0 {
				return
			}

			next, stop := scheduler.GetNextExecution(timeElapsed, executed)
			if stop {
				return
			}
			time.Sleep(next)

			ticks <- struct{}{}
			executed++
		}
	}()

	return results
}

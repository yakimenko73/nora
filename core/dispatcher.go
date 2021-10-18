package core

import (
	"context"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/scheduler"
	"sync"
	"time"
)

func Dispatch(ctx context.Context, scheduler scheduler.Scheduler, executor executor.Executor, duration time.Duration, workers int) <-chan *metric.Result {
	var wg sync.WaitGroup

	ticks := make(chan interface{})
	results := make(chan *metric.Result)
	workerCtx, workerCancel := context.WithCancel(ctx)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			executor.ScheduleExecution(workerCtx, ticks, results)
		}()
	}

	go func() {
		defer close(results)
		defer wg.Wait()
		defer workerCancel()
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
			time.Sleep(next) // FIXME possible deadlock while using with context.WithDeadLine

			ticks <- struct{}{}
			executed++
		}
	}()

	return results
}

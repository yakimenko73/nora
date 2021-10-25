package core

import (
	"context"
	"github.com/illatior/nora/core/executor"
	"github.com/illatior/nora/core/metric"
	"runtime"
	"sync"
	"time"
)

type Dispatcher struct {
	exec executor.Executor

	configuration *LoadOptions
}

func NewDispatcher(opts ...Option) (*Dispatcher, error) {
	d := &Dispatcher{
		configuration: &LoadOptions{
			Duration:  5 * time.Second,
			Workers:   runtime.GOMAXPROCS(0),
			Frequency: 1,
			Period:    1 * time.Second,
		},
		exec: executor.New(),
	}

	var err error
	for _, opt := range opts {
		err = opt.apply(d)
		if err != nil {
			return nil, err
		}
	}

	return d, err
}

func (d *Dispatcher) Dispatch(ctx context.Context) <-chan *metric.Result {
	var wg sync.WaitGroup

	ticks := make(chan interface{}, d.configuration.Workers)
	results := make(chan *metric.Result, d.configuration.Workers)
	ctx, workerCancel := context.WithCancel(ctx)

	for i := 0; i < d.configuration.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.exec.ScheduleExecution(ctx, ticks, results)
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
			if timeElapsed > d.configuration.Duration && d.configuration.Duration != 0 {
				return
			}

			next, stop := d.nextExecution(timeElapsed, executed)
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

func (d *Dispatcher) LoadOptions() LoadOptions {
	return *d.configuration
}

func (d *Dispatcher) nextExecution(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	if d.configuration.Frequency == 0 { // infinite frequency
		return 0, false
	}

	interval := uint64(d.configuration.Period) / d.configuration.Frequency
	expectedHits := uint64(elapsed) / interval
	if expectedHits > hits {
		return 0, false
	}

	next := time.Duration((hits + 1) * interval)
	return next - elapsed, false
}

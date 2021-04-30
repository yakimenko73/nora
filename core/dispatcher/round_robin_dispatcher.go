package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"load-testing/config"
	"load-testing/core/job"
	"load-testing/core/metric"
	"load-testing/core/executor"
	"sync/atomic"
	"time"
)

type roundRobinDispatcher struct {
	errorHandlingDone chan bool

	executor      executor.Executor

	jobs          []job.Job
	currentJob    uint64

	rps uint64

	errChan chan error
}

func NewRoundRobinDispatcher(cfg config.LoadTestConfig, executor executor.Executor) Dispatcher {
	return &roundRobinDispatcher{
		errorHandlingDone: make(chan bool),
		executor:          executor,
		rps:               cfg.RequestsPerSecond,
		errChan:           make(chan error, 0),
		currentJob: 0,
	}
}

func (d *roundRobinDispatcher) Dispatch(ctx context.Context, metricConsumer *metric.MetricConsumer) error {
	g, chctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return d.processErrors(chctx)
	})

	tickInterval := time.Second / time.Duration(d.rps)
	ticker := time.NewTicker(tickInterval)

	g.Go(func() error {
		<- ctx.Done()
		ticker.Stop()

		return nil
	})

	for i := 0; i < 4; i++ {
		g.Go(func() error {
			for {
				select {
				case <- ctx.Done():
					return nil
				case <- chctx.Done():
					return nil
				case <- ticker.C:
					go func() {
						indx, err := d.nextIndex()
						if err != nil {
							panic(err)
						}

						go func(indx int, ctx context.Context) {
							err := d.executor.Execute(d.jobs[indx], ctx, metricConsumer)
							if err != nil {
								d.errChan <- err
							}
						}(indx, ctx)
					}()
				}
			}
		})
	}

	return g.Wait()
}

func (d *roundRobinDispatcher) AddJob(id string, job job.Job) error {
	d.jobs = append(d.jobs, job)

	return nil
}

func (d *roundRobinDispatcher) processErrors(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-d.errChan:
			fmt.Println(err)
		}
	}
}

func (d *roundRobinDispatcher) nextIndex() (int, error) {
	if len(d.jobs) == 0 {
		return 0, errors.New("Workers are not ready yet!")
	}

	return int(atomic.AddUint64(&d.currentJob, uint64(1)) % uint64(len(d.jobs))), nil
}

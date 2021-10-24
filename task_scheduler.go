package task_scheduler

import (
	"context"
	"github.com/illatior/task-scheduler/core"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/cui"
	"golang.org/x/sync/errgroup"
	"sync"
)

type taskScheduler struct {
	c cui.ConsoleUserInterface
	d *core.Dispatcher
}

func New(opts ...Option) (*taskScheduler, error) {
	d, err := core.NewDispatcher()
	if err != nil {
		return nil, err
	}

	ts := &taskScheduler{
		c: nil,
		d: d,
	}

	for _, opt := range opts {
		err := opt.apply(ts)
		if err != nil {
			return nil, err
		}
	}

	return ts, nil
}

func (ts *taskScheduler) Run(ctx context.Context) <-chan *metric.Result {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)

	res := ts.d.Dispatch(ctx)

	// TODO add errs chan with exiting after receiving any error and replace errgroup with it
	userRes := make(chan *metric.Result)
	uiRes := make(chan *metric.Result)
	go func() {
		defer close(userRes)
		defer close(uiRes)
		defer eg.Wait()
		defer cancel()

		dispatchDone := make(chan bool, 1)
		cuiDone := make(chan bool, 1)

		eg.Go(func() error {
			return runMetricRepeater(ctx, userRes, uiRes, res, dispatchDone)
		})

		runCiFunc := ts.getRunCuiFunc(ctx, uiRes, cuiDone, dispatchDone)
		eg.Go(func() error {
			return runCiFunc()
		})

		select {
		case <-ctx.Done():
			return
		case <-cuiDone:
			return
		}
	}()

	return userRes
}

func runMetricRepeater(ctx context.Context,
	userCh, uiCh chan<- *metric.Result,
	resCh <-chan *metric.Result,
	done chan<- bool) error {
	defer func() {
		done <- true
	}()

	// TODO find a better solution to duplicate execution results
	trySendToCh := func(ctx context.Context, m *metric.Result, c chan<- *metric.Result) bool {
		select {
		case <-ctx.Done():
			return false
		case c <- m:
			return true
		}
	}
	for m := range resCh {
		if !trySendToCh(ctx, m, userCh) {
			return nil
		}
		if !trySendToCh(ctx, m, uiCh) {
			return nil
		}
	}
	return nil
}

func (ts *taskScheduler) getRunCuiFunc(ctx context.Context,
	ch <-chan *metric.Result,
	cuiDone chan<- bool,
	dispatchDone <-chan bool) func() error {
	if ts.c != nil {
		return func() error {
			return ts.runCui(ctx, ch, cuiDone)
		}
	}

	return func() error {
		defer func() {
			cuiDone <- true
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-dispatchDone:
				return nil
			case <-ch:
				continue
			}
		}
	}
}

// runCui method is blocking
func (ts *taskScheduler) runCui(ctx context.Context, res <-chan *metric.Result, done chan<- bool) error {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer wg.Wait()
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case m := <-res:
				if m == nil {
					continue
				}

				ts.c.AcceptMetric(m)
			}
		}
	}()

	return ts.c.Run(ctx, done)
}

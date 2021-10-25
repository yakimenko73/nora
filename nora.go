package nora

import (
	"context"
	"github.com/illatior/nora/core"
	"github.com/illatior/nora/core/metric"
	"github.com/illatior/nora/cui"
	"golang.org/x/sync/errgroup"
)

type nora struct {
	c cui.ConsoleUserInterface
	d *core.Dispatcher
}

func NewLoadService(opts ...Option) (*nora, error) {
	d, err := core.NewDispatcher()
	if err != nil {
		return nil, err
	}

	ts := &nora{
		c: cui.NewCuiMock(),
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

func (ts *nora) Run(ctx context.Context) <-chan *metric.Result {
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

		eg.Go(func() error {
			return runMetricRepeater(ctx, userRes, uiRes, res, dispatchDone)
		})
		eg.Go(func() error {
			return ts.c.Run(ctx, uiRes, dispatchDone)
		})

		select {
		case <-ctx.Done():
			return
		case <-ts.c.GetDoneChan():
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
package task_scheduler

import (
	"context"
	"github.com/illatior/task-scheduler/core"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/scheduler"
	"github.com/illatior/task-scheduler/cui"
	"github.com/illatior/task-scheduler/cui/screen"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"runtime"
	"sync"
	"time"
)

type taskScheduler struct {
	duration time.Duration

	sch  scheduler.Scheduler
	exec executor.Executor

	executorsCount int

	withCui  bool
	screens  []cui.Screen
	terminal terminalapi.Terminal
}

func New(opts ...Option) (*taskScheduler, error) {
	sch := scheduler.ConstantScheduler{
		Frequency: 1,
		Period:    1 * time.Second,
	}
	exec := executor.New()

	mainScreen, err := screen.NewMainScreen()
	if err != nil {
		return nil, err
	}

	screens := []cui.Screen{
		mainScreen,
	}

	ts := &taskScheduler{
		duration:       10 * time.Second,
		sch:            sch,
		exec:           exec,
		executorsCount: runtime.GOMAXPROCS(0),
		withCui:        false,
		screens:        screens,
		terminal:       nil,
	}

	for _, opt := range opts {
		err = opt.apply(ts)
		if err != nil {
			return nil, err
		}
	}

	return ts, nil
}

//func (ts *taskScheduler) Run(ctx context.Context) <-chan
//
//func (ts *taskScheduler) RunWithRawResults(ctx context.Context) <-chan *metric.Result {
//
//}

func (ts *taskScheduler) Run(ctx context.Context) <-chan *metric.Result {
	ctx, cancel := context.WithCancel(ctx)

	originalRes := core.Dispatch(ctx, ts.sch, ts.exec, ts.duration, ts.executorsCount)

	userRes := make(chan *metric.Result)
	uiRes := make(chan *metric.Result)
	go func() {
		var wg sync.WaitGroup

		defer close(userRes)
		defer close(uiRes)
		defer wg.Wait()
		defer cancel()

		wg.Add(1)
		go runMetricRepeater(ctx, userRes, uiRes, originalRes)

		wg.Add(1)
		runCiFunc := ts.getRunCuiFunc(ctx, uiRes)
		go runCiFunc()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}()

	return userRes
}

func runMetricRepeater(ctx context.Context,
	userCh, uiCh chan<- *metric.Result,
	resCh <-chan *metric.Result) {

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-resCh:
			userCh <- m
			uiCh <- m
		}
	}
}

func (ts *taskScheduler) getRunCuiFunc(ctx context.Context, ch <-chan *metric.Result) func() {
	if ts.withCui {
		return func () {
			ts.runCui(ctx, ch)
		}
	}

	return func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				continue
			}
		}
	}
}

// runCui method is blocking
func (ts *taskScheduler) runCui(ctx context.Context, res <-chan *metric.Result) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer wg.Done()
	defer cancel()

	ui, err := cui.NewCui(ts.terminal, ts.screens...)
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := ui.Run(ctx)
		if err != nil {
			panic(err) // fixme
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-res:
			ui.AcceptMetric(m)
		}
	}
}

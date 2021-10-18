package task_scheduler

import (
	"context"
	"github.com/illatior/task-scheduler/core"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/scheduler"
	"github.com/illatior/task-scheduler/cui"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"runtime"
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

	screens := []cui.Screen{}

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
		opt.apply(ts)
	}

	return ts, nil
}

func (ts *taskScheduler) Run(ctx context.Context) <-chan *metric.Result {
	originalRes := core.Dispatch(ctx, ts.sch, ts.exec, ts.duration, ts.executorsCount)

	userRes := make(chan *metric.Result)
	uiRes := make(chan *metric.Result)
	go func() {
		defer close(userRes)
		defer close(uiRes)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case m := <-originalRes:
					userRes <- m
					uiRes <- m
				}
			}
		}()

		childCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		if ts.withCui {
			go ts.runCui(childCtx, uiRes)
		} else {
			go func() { // mock
				for {
					select {
					case <-childCtx.Done():
						return
					case <-uiRes:
						continue
					}
				}
			}()
		}
	}()

	return userRes
}

// runCui method is blocking
func (ts *taskScheduler) runCui(ctx context.Context, res <-chan *metric.Result) {
	ui, err := cui.NewCui(ts.terminal)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-res:
			ui.AcceptMetric(m)
		}
	}
}

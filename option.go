package task_scheduler

import (
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/scheduler"
	"github.com/illatior/task-scheduler/cui"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"time"
)

type Option interface {
	apply(ts *taskScheduler)
}

type option func(ts *taskScheduler)

func (o option) apply(ts *taskScheduler) {
	o(ts)
}

func WithDuration(d time.Duration) Option {
	return option(func(ts *taskScheduler) {
		ts.duration = d // FIXME add validation
	})
}

func WithScheduler(sch func() scheduler.Scheduler) Option {
	return option(func(ts *taskScheduler) {
		ts.sch = sch()
	})
}

func WithExecutor(exec func() executor.Executor) Option {
	return option(func(ts *taskScheduler) {
		ts.exec = exec()
	})
}

func WithExecutorsCount(c int) Option {
	return option(func(ts *taskScheduler) {
		ts.executorsCount = c
	})
}

func WithConsoleUserInterface(t terminalapi.Terminal) Option {
	return option(func(ts *taskScheduler) {
		ts.withCui = true
		ts.terminal = t
	})
}

func WithoutDefaultScreens() Option {
	return option(func(ts *taskScheduler) {
		ts.screens = nil
	})
}

func WithCustomScreen(s cui.Screen) Option {
	return option(func(ts *taskScheduler) {
		ts.screens = append(ts.screens, s)
	})
}

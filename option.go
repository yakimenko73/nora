package task_scheduler

import (
	"errors"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/scheduler"
	"github.com/illatior/task-scheduler/cui"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"runtime"
	"time"
)

type Option interface {
	apply(ts *taskScheduler) error
}

type option func(ts *taskScheduler) error

func (o option) apply(ts *taskScheduler) error {
	return o(ts)
}

func WithDuration(d time.Duration) Option {
	return option(func(ts *taskScheduler) error {
		if d < 0 {
			return errors.New("duration of scheduling can't be < 0")
		}

		ts.duration = d
		return nil
	})
}

func WithScheduler(sch func() scheduler.Scheduler) Option {
	return option(func(ts *taskScheduler) error {
		ts.sch = sch()
		return nil
	})
}

func WithExecutor(exec func() executor.Executor) Option {
	return option(func(ts *taskScheduler) error {
		ts.exec = exec()
		return nil
	})
}

func WithExecutorsCount(c int) Option {
	return option(func(ts *taskScheduler) error {
		ts.executorsCount = c
		return nil
	})
}

func createTerminal() (terminalapi.Terminal, error) {
	if runtime.GOOS == "windows" {
		return tcell.New()
	}

	return termbox.New(termbox.ColorMode(terminalapi.ColorMode216))
}

func WithConsoleUserInterface() Option {

	return option(func(ts *taskScheduler) error {
		t, err := createTerminal()
		if err != nil {
			return err
		}

		ts.withCui = true
		ts.terminal = t
		return nil
	})
}

func WithoutDefaultScreens() Option {
	return option(func(ts *taskScheduler) error {
		ts.screens = nil
		return nil
	})
}

func WithCustomScreen(s cui.Screen) Option {
	return option(func(ts *taskScheduler) error {
		ts.screens = append(ts.screens, s)
		return nil
	})
}

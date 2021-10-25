package task_scheduler

import (
	"github.com/illatior/task-scheduler/core"
	"github.com/illatior/task-scheduler/cui"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"runtime"
)

type Option interface {
	apply(ts *taskScheduler) error
}

type option func(ts *taskScheduler) error

func (o option) apply(ts *taskScheduler) error {
	return o(ts)
}

func WithLoadOptions(opts ...core.Option) Option {
	return option(func(ts *taskScheduler) error {
		d, err := core.NewDispatcher(opts...)
		if err != nil {
			return err
		}

		ts.d = d
		return nil
	})
}

func createTerminal() (terminalapi.Terminal, error) {
	if runtime.GOOS == "windows" {
		return tcell.New()
	}

	return termbox.New(termbox.ColorMode(terminalapi.ColorMode216))
}

func WithConsoleUserInterface(opts ...cui.Option) Option {

	return option(func(ts *taskScheduler) error {
		t, err := createTerminal()
		if err != nil {
			return err
		}

		ts.c, err = cui.NewCui(t, opts...)
		if err != nil {
			return err
		}

		return nil
	})
}

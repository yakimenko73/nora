package nora

import (
	"github.com/illatior/nora/core"
	"github.com/illatior/nora/cui"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"runtime"
)

type Option interface {
	apply(n *nora) error
}

type option func(n *nora) error

func (o option) apply(n *nora) error {
	return o(n)
}

func WithLoadOptions(opts ...core.Option) Option {
	return option(func(n *nora) error {
		d, err := core.NewDispatcher(opts...)
		if err != nil {
			return err
		}

		n.d = d
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

	return option(func(n *nora) error {
		t, err := createTerminal()
		if err != nil {
			return err
		}

		n.c, err = cui.NewCui(t, opts...)
		if err != nil {
			return err
		}

		return nil
	})
}

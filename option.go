package nora

import (
	"github.com/illatior/nora/cli"
	"github.com/illatior/nora/lib"
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

func WithLoadOptions(opts ...lib.Option) Option {
	return option(func(n *nora) error {
		d, err := lib.NewDispatcher(opts...)
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

func WithConsoleUserInterface(opts ...cli.Option) Option {

	return option(func(n *nora) error {
		t, err := createTerminal()
		if err != nil {
			return err
		}

		n.c, err = cli.NewCui(t, opts...)
		if err != nil {
			return err
		}

		return nil
	})
}

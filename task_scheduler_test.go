package task_scheduler

import (
	"context"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"time"

	"testing"
)

func TestSandbox(t *testing.T) {
	ter, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		t.Fatal(err)
	}

	tsch, err := New(
		WithDuration(time.Second),
		WithConsoleUserInterface(ter),
	)
	if err != nil {
		t.Fatal(err)
	}

	for range tsch.Run(context.Background()) {
	}
}
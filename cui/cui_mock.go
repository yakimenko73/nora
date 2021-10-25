package cui

import (
	"context"
	"github.com/illatior/nora/core/metric"
)

type cuiMock struct {
	done chan bool
}

func NewCuiMock() ConsoleUserInterface {
	return &cuiMock{
		done: make(chan bool),
	}
}

func (ui *cuiMock) Run(ctx context.Context, metrics <-chan *metric.Result, dispatchDone <-chan bool) error {
	defer func() {
		ui.done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-dispatchDone:
			return nil
		case <-metrics:
			continue
		}
	}
}

func (ui *cuiMock) GetDoneChan() <-chan bool {
	return ui.done
}

func (ui *cuiMock) ChangeFullscreenState() error {
	return nil
}

func (ui *cuiMock) NextScreen() error {
	return nil
}

func (ui *cuiMock) PreviousScreen() error {
	return nil
}

func (ui *cuiMock) IsFullscreen() bool {
	return true
}

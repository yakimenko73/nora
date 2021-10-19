package cui

import (
	"context"
	"errors"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"sync"
)

type cui struct {
	isFullscreen bool

	t terminalapi.Terminal
	c *container.Container

	mu            sync.RWMutex
	currentScreen int
	screens       []Screen
}

func NewCui(t terminalapi.Terminal, screens ...Screen) (ConsoleUserInterface, error) {
	c, err := container.New(
		t,
		container.ID(SCREEN_ID),
		container.Border(linestyle.Light),
		container.BorderTitle("Task-scheduler"),
	)
	if err != nil {
		return nil, err
	}

	ui := &cui{
		isFullscreen: false,
		c:            c,
		t:            t,

		mu:            sync.RWMutex{},
		currentScreen: 0,
		screens:       screens,
	}
	err = ui.changeMainScreen()
	if err != nil {
		return nil, err
	}

	return ui, nil
}

func (ui *cui) Run(ctx context.Context) error {
	childCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		<- ctx.Done()
	}()

	subs := func (k *terminalapi.Keyboard) {
		var err error
		switch k.Key {
		case 'Q', 'q', keyboard.KeyCtrlC:
			cancel()
		case 'A', 'a':
			err = ui.PreviousScreen()
		case 'D', 'd':
			err = ui.NextScreen()
		case 'F', 'f':
			err = ui.ChangeFullscreenState()
		default:
			return
		}

		if err != nil {
			panic(err)
		}
	}

	// TODO add keyboard subscriber's
	return termdash.Run(childCtx, ui.t, ui.c, termdash.KeyboardSubscriber(subs))
}

func (ui *cui) AcceptMetric(m *metric.Result) {
	panic(errors.New("qwe"))
}

func (ui *cui) ChangeFullscreenState() error {
	ui.isFullscreen = !ui.isFullscreen

	return ui.changeMainScreen()
}

func (ui *cui) NextScreen() error {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.currentScreen++
	if ui.currentScreen == len(ui.screens)-1 {
		ui.currentScreen = 0
	}

	return ui.changeMainScreen()
}

func (ui *cui) PreviousScreen() error {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.currentScreen--
	if ui.currentScreen < 0 {
		ui.currentScreen = len(ui.screens) - 1
	}

	return ui.changeMainScreen()
}

func (ui *cui) changeMainScreen() error {
	currentScreen := ui.screens[ui.currentScreen]

	builder := grid.New()
	body := currentScreen.GetBody()
	if ui.isFullscreen {
		builder.Add(body)
	} else {
		builder.Add(
			currentScreen.GetHeader(),
			currentScreen.GetBody(),
			currentScreen.GetFooter(),
		)
	}

	opts, err := builder.Build()
	if err != nil {
		return err
	}

	return ui.c.Update(SCREEN_ID, opts...)
}

func (ui *cui) IsFullscreen() bool {
	return ui.isFullscreen
}

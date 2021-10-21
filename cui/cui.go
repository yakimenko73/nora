package cui

import (
	"context"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/cui/screen"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"sync"
	"time"
)

type subsFunc func (ctx context.Context, cancel context.CancelFunc, ui *cui) func(*terminalapi.Keyboard)

type cui struct {
	isFullscreen bool

	t terminalapi.Terminal
	c *container.Container

	screenMu      sync.RWMutex
	currentScreen int
	screens       []screen.Screen

	displayInterval time.Duration
	updateInterval  time.Duration

	metricsMu sync.RWMutex
	metrics   metric.Metrics

	subs subsFunc
}

func NewCui(t terminalapi.Terminal, opts ...Option) (ConsoleUserInterface, error) {
	c, err := container.New(
		t,
		container.ID(SCREEN_ID),
		container.Border(linestyle.Light),
		container.BorderTitle("Task-scheduler"),
	)
	if err != nil {
		return nil, err
	}

	mainScreen, err := screen.NewMainScreen()
	if err != nil {
		return nil, err
	}

	ui := &cui{
		isFullscreen: false,
		c:            c,
		t:            t,

		screenMu:      sync.RWMutex{},
		currentScreen: 0,
		screens:       []screen.Screen{mainScreen},

		displayInterval: 60 * time.Second,
		updateInterval:  100 * time.Millisecond,

		metrics: metric.NewMetrics(),

		subs: defaultSubs(),
	}

	for _, opt := range opts {
		err = opt.apply(ui)
		if err != nil {
			return nil, err
		}
	}

	err = ui.changeMainScreen()
	if err != nil {
		return nil, err
	}

	return ui, nil
}

func (ui *cui) Run(ctx context.Context, done chan<- bool) error {
	defer func() {
		done <- true
	}()

	ctx, cancel := context.WithCancel(ctx)
	go ui.update(ctx)
	go func() {
		for _, s := range ui.screens {
			go s.Run(ctx)
		}
	}()

	// TODO add ability to customize subs with option
	// TODO add +/- handlers for increasing/decreasing displayable period
	subs := ui.subs(ctx, cancel, ui)

	defer func() {
		ui.t.Close()
	}()
	return termdash.Run(ctx, ui.t, ui.c, termdash.KeyboardSubscriber(subs))
}

func (ui *cui) AcceptMetric(m *metric.Result) {
	ui.metricsMu.Lock()
	defer ui.metricsMu.Unlock()

	ui.metrics.ConsumeResult(m)
}

func (ui *cui) update(ctx context.Context) {
	t := time.NewTicker(ui.updateInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			currentScreen := ui.screens[ui.currentScreen]

			currentScreen.GetMetricsChan() <- ui.metrics
		}
	}
}

func (ui *cui) ChangeFullscreenState() error {
	ui.isFullscreen = !ui.isFullscreen

	return ui.changeMainScreen()
}

func (ui *cui) NextScreen() error {
	ui.screenMu.Lock()
	defer ui.screenMu.Unlock()

	ui.currentScreen++
	if ui.currentScreen == len(ui.screens) {
		ui.currentScreen = 0
	}

	return ui.changeMainScreen()
}

func (ui *cui) PreviousScreen() error {
	ui.screenMu.Lock()
	defer ui.screenMu.Unlock()

	ui.currentScreen--
	if ui.currentScreen < 0 {
		ui.currentScreen = len(ui.screens) - 1
	}

	return ui.changeMainScreen()
}

func (ui *cui) changeMainScreen() error { // FIXME after exiting fullscreen mode main BorderTitle and BorderStyle continues to be as body's one
	currentScreen := ui.screens[ui.currentScreen]

	builder := grid.New()
	if ui.isFullscreen {
		addElem(currentScreen.GetBody(), builder)
	} else {
		addElem(currentScreen.GetHeader(), builder)
		addElem(currentScreen.GetBody(), builder)
		addElem(currentScreen.GetFooter(), builder)
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

func defaultSubs() subsFunc {
	return func (ctx context.Context, cancel context.CancelFunc, ui *cui) func(*terminalapi.Keyboard) {
		return func(k *terminalapi.Keyboard) {
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
	}
}

func addElem(e grid.Element, b *grid.Builder) {
	if e != nil {
		b.Add(e)
	}
}

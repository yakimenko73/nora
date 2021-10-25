package cui

import (
	"context"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

type subsFunc func(ctx context.Context, cancel context.CancelFunc, ui *cui) func(*terminalapi.Keyboard)

func defaultSubs() subsFunc {
	return func(ctx context.Context, cancel context.CancelFunc, ui *cui) func(*terminalapi.Keyboard) {
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
			case '+':
				ui.screens[ui.currentScreen].ChangeDisplayInterval(ui.changeDisplayableIntervalDelta)
			case '-':
				ui.screens[ui.currentScreen].ChangeDisplayInterval(-ui.changeDisplayableIntervalDelta)
			default:
				return
			}

			if err != nil {
				panic(err)
			}
		}
	}
}

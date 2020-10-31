package gui

import (
	"context"

	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func keybinds(cancel context.CancelFunc) func(*terminalapi.Keyboard) {
	return func(k *terminalapi.Keyboard) {
		switch k.Key {
		case keyboard.KeyCtrlC, 'q': // Quit
			cancel()
		}
	}
}

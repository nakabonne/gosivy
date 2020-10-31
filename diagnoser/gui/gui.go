package gui

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/nakabonne/gosivy/stats"
)

const (
	// How often termdash redraws the screen.
	defaultRedrawInterval = time.Second
	rootID                = "root"
)

type runner func(ctx context.Context, t terminalapi.Terminal, c *container.Container, opts ...termdash.Option) error

// Run stats drawing charts, and blocks until the quit operation is performed.
func Run(ctx context.Context, redrawIntarval time.Duration, meta *stats.Meta, statsCh <-chan *stats.Stats) error {
	var (
		t   terminalapi.Terminal
		err error
	)
	if runtime.GOOS == "windows" {
		t, err = tcell.New()
	} else {
		t, err = termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	}
	if err != nil {
		return fmt.Errorf("failed to generate terminal interface: %w", err)
	}
	defer t.Close()
	return run(ctx, t, termdash.Run)
}

func run(ctx context.Context, t terminalapi.Terminal, r runner) error {
	return nil
}

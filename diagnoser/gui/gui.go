// Package gui provides an ability to draw charts on the terminal.
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
	defaultRedrawInterval = time.Second
	rootID                = "root"
)

type GUI struct {
	// How often termdash redraws the screen.
	RedrawInterval time.Duration
	// The function to quit the application.
	Cancel context.CancelFunc
	// A channel for receiving data sources to draw on the chart.
	StatsCh <-chan *stats.Stats
	// Metadata of the process where the agent runs on.
	Metadata stats.Meta
}

func NewGUI(redrawInterval time.Duration, cancel context.CancelFunc, statsCh <-chan *stats.Stats, metadata *stats.Meta) *GUI {
	if redrawInterval == 0 {
		redrawInterval = defaultRedrawInterval
	}
	if statsCh == nil {
		statsCh = make(<-chan *stats.Stats)
	}
	return &GUI{
		RedrawInterval: redrawInterval,
		Cancel:         cancel,
		StatsCh:        statsCh,
		Metadata:       *metadata,
	}
}

// Run starts drawing charts, and blocks until the quit operation is performed.
func (g *GUI) Run(ctx context.Context) error {
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
	return g.run(ctx, t, termdash.Run)
}

type runner func(ctx context.Context, t terminalapi.Terminal, c *container.Container, opts ...termdash.Option) error

func (g *GUI) run(ctx context.Context, t terminalapi.Terminal, r runner) error {
	c, err := container.New(t, container.ID(rootID))
	if err != nil {
		return fmt.Errorf("failed to generate container: %w", err)
	}

	w, err := newWidgets()
	if err != nil {
		return fmt.Errorf("failed to generate widgets: %w", err)
	}

	gridOpts, err := gridLayout(w)
	if err != nil {
		return fmt.Errorf("failed to build grid layout: %w", err)
	}

	if err := c.Update(rootID, gridOpts.base...); err != nil {
		return fmt.Errorf("failed to update container: %w", err)
	}
	k := keybinds(ctx, g.Cancel)

	return r(ctx, t, c, termdash.KeyboardSubscriber(k), termdash.RedrawInterval(g.RedrawInterval))
}

// gridOpts holds all options in our grid. It basically holds the container
// options (column, width, padding, etc) of our widgets.
type gridOpts struct {
	base []container.Option
}

func gridLayout(w *widgets) (*gridOpts, error) {
	return nil, nil
}

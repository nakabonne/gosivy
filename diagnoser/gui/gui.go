// Package gui provides an ability to draw charts on the terminal.
package gui

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"

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

	widgets *widgets
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

	g.widgets, err = newWidgets(&g.Metadata)
	if err != nil {
		return fmt.Errorf("failed to generate widgets: %w", err)
	}

	opts, err := gridLayout(g.widgets)
	if err != nil {
		return fmt.Errorf("failed to build grid layout: %w", err)
	}

	if err := c.Update(rootID, opts...); err != nil {
		return fmt.Errorf("failed to update container: %w", err)
	}

	go g.appendStats(ctx)

	k := keybinds(g.Cancel)

	return r(ctx, t, c, termdash.KeyboardSubscriber(k), termdash.RedrawInterval(g.RedrawInterval))
}

// gridLayout gives back options for grid layout, which is composed by rows that inside of the rows
// are the columns, the elements are the columns.
//
// ----------------------------------------------------
// [------element------] [--element--] [---element---]
// ----------------------------------------------------
// [element] [element] [------------element----------]
// ----------------------------------------------------
// [-element-]       [----element----]        [element]
// ----------------------------------------------------
func gridLayout(w *widgets) ([]container.Option, error) {
	raw1 := grid.RowHeightPerc(4,
		grid.Widget(w.Metadata, container.Border(linestyle.Light), container.BorderTitle("Press Q to quit")),
	)
	raw2 := grid.RowHeightPerc(45,
		grid.ColWidthPerc(50, grid.Widget(w.CPUChart, container.Border(linestyle.Light), container.BorderTitle("CPU Usage (%)"))),
		grid.ColWidthPerc(50, grid.Widget(w.GoroutineChart, container.Border(linestyle.Light), container.BorderTitle("Goroutines"))),
	)
	raw3 := grid.RowHeightPercWithOpts(45,
		[]container.Option{container.Border(linestyle.Light), container.BorderTitle("Heap (MB)")},
		grid.RowHeightPerc(97, grid.ColWidthPerc(99, grid.Widget(w.HeapChart))),
		grid.RowHeightPercWithOpts(3,
			[]container.Option{container.MarginLeftPercent(0), container.MarginBottomPercent(0)},
			textsInColumn(w.HeapAllocLegend.text, w.HeapIdelLegend.text, w.HeapInuseLegend.text)...,
		),
	)
	builder := grid.New()
	builder.Add(
		raw1,
		raw2,
		raw3,
	)

	return builder.Build()
}

func textsInColumn(texts ...Text) []grid.Element {
	els := make([]grid.Element, 0, len(texts))
	for _, text := range texts {
		els = append(els, grid.ColWidthPerc(3, grid.Widget(text)))
	}
	return els
}

// appendStats appends entities as soon as a stats arrives.
// Note that it doesn't redraw the moment stats are appended.
func (g *GUI) appendStats(ctx context.Context) {
	var (
		cpuUsages  = make([]float64, 0)
		goroutines = make([]float64, 0)
		allocs     = make([]float64, 0)
		idles      = make([]float64, 0)
		inuses     = make([]float64, 0)

		megaByte uint64 = 1000000
	)

	for {
		select {
		case <-ctx.Done():
			return
		case stats := <-g.StatsCh:
			if stats == nil {
				continue
			}
			cpuUsages = append(cpuUsages, stats.CPUUsage)
			goroutines = append(goroutines, float64(stats.Goroutines))
			allocs = append(allocs, float64(stats.HeapAlloc/megaByte))
			idles = append(idles, float64(stats.HeapIdle/megaByte))
			inuses = append(inuses, float64(stats.HeapInuse/megaByte))

			g.widgets.CPUChart.Series("cpu-usage", cpuUsages,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(87))),
			)
			g.widgets.GoroutineChart.Series("goroutines", goroutines,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(87))),
			)
			g.widgets.HeapChart.Series("alloc", allocs,
				linechart.SeriesCellOpts(g.widgets.HeapAllocLegend.cellOpts...),
			)
			g.widgets.HeapChart.Series("idle", idles,
				linechart.SeriesCellOpts(g.widgets.HeapIdelLegend.cellOpts...),
			)
			g.widgets.HeapChart.Series("inuse", inuses,
				linechart.SeriesCellOpts(g.widgets.HeapInuseLegend.cellOpts...),
			)
		}
	}
}

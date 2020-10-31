package gui

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"

	"github.com/nakabonne/gosivy/stats"
)

type LineChart interface {
	widgetapi.Widget
	Series(label string, values []float64, opts ...linechart.SeriesOption) error
}

type Text interface {
	widgetapi.Widget
	Write(text string, wOpts ...text.WriteOption) error
}

type chartLegend struct {
	text     Text
	cellOpts []cell.Option
}

type widgets struct {
	Metadata Text

	CPUChart       LineChart
	GoroutineChart LineChart
	HeapChart      LineChart

	HeapAllocLegend chartLegend
	HeapIdelLegend  chartLegend
	HeapInuseLegend chartLegend
}

func newWidgets(meta *stats.Meta) (*widgets, error) {
	metadata, err := newText(meta.String())
	if err != nil {
		return nil, err
	}

	cpuChart, err := newLineChart()
	if err != nil {
		return nil, err
	}
	goroutineChart, err := newLineChart()
	if err != nil {
		return nil, err
	}
	heapChart, err := newLineChart()
	if err != nil {
		return nil, err
	}

	allocColor := cell.FgColor(cell.ColorMagenta)
	allocText, err := newText("alloc", text.WriteCellOpts(allocColor))
	if err != nil {
		return nil, err
	}
	idleColor := cell.FgColor(cell.ColorGreen)
	idleText, err := newText("idle", text.WriteCellOpts(idleColor))
	if err != nil {
		return nil, err
	}
	inuseColor := cell.FgColor(cell.ColorYellow)
	inuseText, err := newText("idle", text.WriteCellOpts(inuseColor))
	if err != nil {
		return nil, err
	}
	return &widgets{
		Metadata:        metadata,
		CPUChart:        cpuChart,
		GoroutineChart:  goroutineChart,
		HeapChart:       heapChart,
		HeapAllocLegend: chartLegend{allocText, []cell.Option{allocColor}},
		HeapIdelLegend:  chartLegend{idleText, []cell.Option{idleColor}},
		HeapInuseLegend: chartLegend{inuseText, []cell.Option{inuseColor}},
	}, nil
}

func newLineChart() (LineChart, error) {
	return linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
}

func newText(s string, opts ...text.WriteOption) (Text, error) {
	t, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		return nil, err
	}
	if s != "" {
		if err := t.Write(s, opts...); err != nil {
			return nil, err
		}
	}
	return t, nil
}

package screen

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

type screenOpts struct {
	header grid.Element
	body   grid.Element
	footer grid.Element
}

func newLineChart() (*linechart.LineChart, error) {
	return linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
}

func newTextBlock() (*text.Text, error) {
	return text.New(text.WrapAtWords(), text.RollContent())
}

func newBarChart() (*barchart.BarChart, error) {
	return barchart.New() // TODO
}

func borderLight() container.Option {
	return container.Border(linestyle.Light)
}

func borderTitle(t string) container.Option {
	return container.BorderTitle(t)
}

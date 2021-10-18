package screen

import (
	"github.com/illatior/task-scheduler/cui"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

type mainScreen struct {
	opts screenOpts

	// header
	// `empty`

	// body
	latencyChart *linechart.LineChart

	// footer
	optionsText   *text.Text
	latenciesText *text.Text
	responsesText *text.Text
	errorsText    *text.Text
	errorsCount   *barchart.BarChart

	// todo add information for debug like duration, requests, responses, throughput, ...
}

func NewMainScreen() (cui.Screen, error) {
	m, err := buildMainScreen()
	if err != nil {
		return nil, err
	}

	headerOpts := grid.RowHeightPerc(0)

	body := grid.RowHeightPercWithOpts(
		70,
		[]container.Option{container.ID(cui.SCREEN_ID)},
		grid.Widget(m.latencyChart, borderLight(), borderTitle("Latency (ms)")),
	)

	footerOpts := grid.RowHeightPerc(
		30,
		grid.ColWidthPerc(10, grid.Widget(m.optionsText, borderLight(), borderTitle("Run options"))),
		grid.ColWidthPerc(20, grid.Widget(m.latenciesText, borderLight(), borderTitle("Latencies"))),
		grid.ColWidthPerc(20, grid.Widget(m.responsesText, borderLight(), borderTitle("Responses"))),
		grid.ColWidthPerc(20, grid.Widget(m.errorsText, borderLight(), borderTitle("Errors"))),
		grid.ColWidthPerc(30, grid.Widget(m.errorsCount, borderLight(), borderTitle("Errors count"))),
	)

	m.opts = screenOpts{
		header: headerOpts,
		body:   body,
		footer: footerOpts,
	}

	return m, nil
}

func buildMainScreen() (*mainScreen, error) {
	latencyChart, err := newLineChart()
	if err != nil {
		return nil, err
	}

	optionsText, err := newTextBlock()
	if err != nil {
		return nil, err
	}
	latenciesText, err := newTextBlock()
	if err != nil {
		return nil, err
	}
	responsesText, err := newTextBlock()
	if err != nil {
		return nil, err
	}
	errorsText, err := newTextBlock()
	if err != nil {
		return nil, err
	}
	errorsCount, err := newBarChart()
	if err != nil {
		return nil, err
	}

	return &mainScreen{
		latencyChart:  latencyChart,
		optionsText:   optionsText,
		latenciesText: latenciesText,
		responsesText: responsesText,
		errorsText:    errorsText,
		errorsCount:   errorsCount,
	}, nil
}

func (s *mainScreen) GetBody() grid.Element {
	return s.opts.body
}

func (s *mainScreen) GetHeader() grid.Element {
	return s.opts.header
}

func (s *mainScreen) GetFooter() grid.Element {
	return s.opts.footer
}

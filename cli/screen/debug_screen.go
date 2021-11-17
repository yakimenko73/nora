package screen

import (
	"github.com/illatior/nora/lib/metric"
	"github.com/mum4k/termdash/container/grid"
)

type debugScreen struct {
}

func NewDebugScreen() (Screen, error) {
	return nil, nil
}

func (s *debugScreen) GetBody() grid.Element {
	return nil
}

func (s *debugScreen) GetHeader() grid.Element {
	return nil
}

func (s *debugScreen) GetFooter() grid.Element {
	return nil
}

func (s *debugScreen) UpdateWithLatencyMetrics(m metric.LatencyMetrics) {

}

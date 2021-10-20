package screen

import (
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/cui"
	"github.com/mum4k/termdash/container/grid"
)

type debugScreen struct {
}

func NewDebugScreen() (cui.Screen, error) {
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

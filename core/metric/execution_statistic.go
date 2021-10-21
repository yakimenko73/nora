package metric

import (
	"sync"
	"sync/atomic"
)

type executionStatistic struct {
	totalExecuted uint64
	totalSuccess  uint64

	mu sync.Mutex
	errors map[string]uint64
}

func NewExecutionStatistic() ExecutionStatistic {
	return &executionStatistic{
		totalExecuted: 0,
		totalSuccess:  0,
		errors:        make(map[string]uint64),
	}
}

func (es *executionStatistic) ConsumeResult(res *Result) {
	atomic.AddUint64(&es.totalExecuted, 1)
	if res.Error != nil {
		if _, ok := es.errors[res.Error.Error()]; !ok {
			es.errors[res.Error.Error()] = 0
		}

		es.mu.Lock()
		defer es.mu.Unlock()
		es.errors[res.Error.Error()] = es.errors[res.Error.Error()] + 1 // TODO find a better solution

		return
	}

	atomic.AddUint64(&es.totalSuccess, 1)
}

func (es *executionStatistic) GetTotalExecuted() uint64 {
	return es.totalExecuted
}

func (es *executionStatistic) GetTotalSuccess() uint64 {
	return es.totalSuccess
}

func (es *executionStatistic) GetErrors() map[string]uint64 {
	return es.errors
}

package metric

import (
	"sync"
	"time"
)

// TODO add OrderStaticsTree to get O(logn) time complexity of getting any percentile and do not calculate percentiles for every new execution result
type latencyMetrics struct {
	mu sync.Mutex

	total uint64

	min time.Duration
	max time.Duration
	avg float64
}

func NewLatencyMetrics() LatencyMetrics {
	return &latencyMetrics{
		mu:    sync.Mutex{},
		total: 0,
		min:   -1,
		max:   -1,
		avg:   0,
	}
}

func (lm *latencyMetrics) ConsumeResult(res *Result) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.total++
	if lm.min > res.Duration || lm.min == -1 {
		lm.min = res.Duration
	}
	if lm.max < res.Duration || lm.max == -1 {
		lm.max = res.Duration
	}

	prevSum := lm.avg*float64(lm.total-1)
	curSum := prevSum + float64(res.Duration)
	lm.avg = curSum / float64(lm.total)


}

func (lm *latencyMetrics) GetPercentile(p int) time.Duration {
	return 0
}

func (lm *latencyMetrics) GetMin() time.Duration {
	return lm.min
}

func (lm *latencyMetrics) GetMax() time.Duration {
	return lm.max
}

func (lm *latencyMetrics) GetAvg() time.Duration {
	return time.Duration(uint64(lm.avg))
}

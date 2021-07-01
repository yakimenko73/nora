package metric

import (
	"sort"
	"time"
)

type Metrics struct {
	Results map[string][]*Result
}

func (m *Metrics) ConsumeResult(result *Result) {
	if m.Results == nil {
		m.Results = make(map[string][]*Result, 0)
	}

	if _, ok := m.Results[result.Name]; !ok {
		m.Results[result.Name] = make([]*Result, 0)
	}

	m.Results[result.Name] = append(m.Results[result.Name], result)
}

func (m *Metrics) GetLatencyMetrics() map[string]LatencyMetrics  {
	latencies := make(map[string]LatencyMetrics, 0)
	for job, metrics := range m.Results {

		durations := make([]time.Duration, 0)
		lMetrics := LatencyMetrics{
			Min: -1,
			Max: -1,
			Total: 0,
		}
		for _, metric := range metrics {
			lMetrics.Total += metric.Duration
			if lMetrics.Min > metric.Duration || lMetrics.Min == -1 {
				lMetrics.Min = metric.Duration
			}
			if lMetrics.Max < metric.Duration || lMetrics.Max == -1 {
				lMetrics.Max = metric.Duration
			}

			durations = SortedInsert(durations, metric.Duration)
		}
		lMetrics.Q1 = calculatePercentile(durations, 25)
		lMetrics.Median = calculatePercentile(durations, 50)
		lMetrics.Q3 = calculatePercentile(durations, 75)
		lMetrics.P90 = calculatePercentile(durations, 90)
		lMetrics.P95 = calculatePercentile(durations, 95)
		lMetrics.P99 = calculatePercentile(durations, 99)

		latencies[job] = lMetrics
	}

	return latencies
}

func SortedInsert(durations []time.Duration, duration time.Duration) []time.Duration {
	i := sort.Search(len(durations), func(i int) bool {return durations[i] >= duration})
	durations = append(durations, 0)
	copy(durations[i+1:],durations[i:])
	durations[i] = duration
	return durations
}
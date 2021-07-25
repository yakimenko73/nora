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

func (m *Metrics) GetExecutionStatistic() ExecutionStatistic {
	em := ExecutionStatistic{
		TotalRequests:     0,
		RequestsPerSecond: 0,
		TotalSuccess:      0,
		ErrorCodes:        make(map[string]int64, 0),
	}

	for _, metrics := range m.Results {
		em.TotalRequests += int64(len(metrics))

		for _, metric := range metrics {
			if metric.Start.Before(em.StartTime) || em.StartTime.Equal(time.Time{}) {
				em.StartTime = metric.Start
			}
			if metricEndTime := metric.Start.Add(metric.Duration); metricEndTime.After(em.EndTime) || em.EndTime.Equal(time.Time{}) {
				em.EndTime = metricEndTime
			}

			if metric.Error != nil {
				if _, ok := em.ErrorCodes[metric.Error.Error()]; !ok {
					em.ErrorCodes[metric.Error.Error()] = 0
				}

				em.ErrorCodes[metric.Error.Error()]++
				continue
			}

			em.TotalSuccess += int64(1)
		}
	}

	em.Duration = em.EndTime.Sub(em.StartTime)
	em.RequestsPerSecond = em.TotalRequests / int64(em.Duration/time.Second)

	return em
}

func (m *Metrics) GetLatencyMetrics() map[string]LatencyMetrics {
	latencies := make(map[string]LatencyMetrics, 0)
	for task, metrics := range m.Results {

		durations := make([]time.Duration, 0)
		lMetrics := LatencyMetrics{
			Min:   -1,
			Max:   -1,
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

			durations = sortedInsert(durations, metric.Duration)
		}
		lMetrics.Q1 = calculatePercentile(durations, 25)
		lMetrics.Median = calculatePercentile(durations, 50)
		lMetrics.Q3 = calculatePercentile(durations, 75)
		lMetrics.P90 = calculatePercentile(durations, 90)
		lMetrics.P95 = calculatePercentile(durations, 95)
		lMetrics.P99 = calculatePercentile(durations, 99)

		latencies[task] = lMetrics
	}

	return latencies
}

func (m *Metrics) GetLineChartsReport(groupingFactor time.Duration) *LineChartsReport {
	report := &LineChartsReport{
		Charts: make(map[string][]*ChartData, 0),
	}

	for job, metrics := range m.Results {
		if _, ok := report.Charts[job]; !ok {
			report.Charts[job] = make([]*ChartData, 0)
		}

		for _, metric := range metrics {
			cd := &ChartData{
				Duration: metric.Duration,
				Time:     metric.Start.Truncate(groupingFactor),
			}
			report.Charts[job] = append(report.Charts[job], cd)
		}
	}

	return report
}

func sortedInsert(durations []time.Duration, duration time.Duration) []time.Duration {
	i := sort.Search(len(durations), func(i int) bool { return durations[i] >= duration })
	durations = append(durations, 0)
	copy(durations[i+1:], durations[i:])
	durations[i] = duration
	return durations
}

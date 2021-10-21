package metric

import "time"

type chartMetrics struct {
	entries []*ChartEntry
}

func NewChartMetrics() ChartMetrics {
	return &chartMetrics{
		entries: make([]*ChartEntry, 0),
	}
}

func (cm *chartMetrics) ConsumeResult(res *Result) {
	entry := &ChartEntry{
		Timestamp: res.Start, // TODO prob res.Start.Add(res.Duration) will be more correct
		Duration:  res.Duration,
	}

	cm.entries = sortedInsert(cm.entries, entry)
}

func (cm *chartMetrics) GetInRange(from, to time.Time) []ChartEntry {
	res := make([]ChartEntry, 0)
	for i := len(cm.entries)-1; i >= 0; i-- {
		entry := *cm.entries[i]
		if !(entry.Timestamp.After(to) || entry.Timestamp.Before(from)) {
			res = append(res, entry)
		}
	}

	return res
}
package metric

import "time"

type Result struct {
	Name string

	Start    time.Time
	Duration time.Duration

	Error error
}

type LatencyMetrics struct {
	Total time.Duration

	Min time.Duration
	Max time.Duration

	Q1     time.Duration
	Median time.Duration
	Q3     time.Duration

	P90 time.Duration
	P95 time.Duration
	P99 time.Duration
}

type ExecutionStatistic struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration

	TotalRequests     int64
	RequestsPerSecond int64

	TotalSuccess int64
	ErrorCodes   map[string]int64
}

type LineChartsReport struct {
	Charts map[string][]*ChartData
}

type ChartData struct {
	Duration time.Duration
	Time     time.Time
}

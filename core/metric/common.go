package metric

import "time"

type Result struct {
	Name string

	Start    time.Time
	Duration time.Duration

	Error           error
}

type LatencyMetrics struct {
	Total time.Duration

	Min time.Duration
	Max time.Duration

	Q1 time.Duration
	Q3 time.Duration
	Median time.Duration

	P90 time.Duration
	P95 time.Duration
	P99 time.Duration
}
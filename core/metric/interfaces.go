package metric

import "time"

type Metrics interface {
	ConsumeResult(res *Result)

	GetLatencyMetrics() LatencyMetrics
	GetExecutionStatistic() ExecutionStatistic
	GetChartMetrics() ChartMetrics

	String() string
}

type LatencyMetrics interface {
	ConsumeResult(res *Result)

	GetPercentile(p int) time.Duration
	GetMin() time.Duration
	GetMax() time.Duration
	GetAvg() time.Duration

	String() string
}

type ExecutionStatistic interface {
	ConsumeResult(res *Result)

	GetTotalExecuted() uint64
	GetTotalSuccess() uint64
	GetErrors() []string
	GetErrorsCount(err string) uint64

	String() string
}

type ChartMetrics interface {
	ConsumeResult(res *Result)
	GetInRange(from, to time.Time) []ChartEntry

	String() string
}

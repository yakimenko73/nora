package metric

type metrics struct {
	latencyMetrics     LatencyMetrics
	executionStatistic ExecutionStatistic
	chartMetrics       ChartMetrics
}

func NewMetrics() Metrics {
	return &metrics{
		latencyMetrics:     NewLatencyMetrics(),
		executionStatistic: NewExecutionStatistic(),
		chartMetrics:       NewChartMetrics(),
	}
}

func (m *metrics) ConsumeResult(result *Result) {
	m.latencyMetrics.ConsumeResult(result)
	m.executionStatistic.ConsumeResult(result)
	m.chartMetrics.ConsumeResult(result)
}

func (m *metrics) GetLatencyMetrics() LatencyMetrics {
	return m.latencyMetrics
}

func (m *metrics) GetExecutionStatistic() ExecutionStatistic {
	return m.executionStatistic
}

func (m *metrics) GetChartMetrics() ChartMetrics {
	return m.chartMetrics
}

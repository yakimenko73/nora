package metric

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetLatencyMetrics_SingleResult(t *testing.T) {
	t.Parallel()

	metricName := "qwe"

	m := Metrics{}
	m.ConsumeResult(
		&Result{
			Name:     metricName,
			Start:    time.Now(),
			Duration: time.Second,
			Error:    nil,
		},
	)

	res := m.GetLatencyMetrics()
	assert.Equal(t, len(res), 1)

	actualMetric := res[metricName]
	assert.NotNil(t, actualMetric)

	assert.Equal(t, actualMetric.Max, time.Second)
	assert.Equal(t, actualMetric.Min, time.Second)
	assert.Equal(t, actualMetric.Total, time.Second)
	assert.Equal(t, actualMetric.Q1, time.Second)
	assert.Equal(t, actualMetric.Q3, time.Second)
	assert.Equal(t, actualMetric.Median, time.Second)
	assert.Equal(t, actualMetric.P90, time.Second)
	assert.Equal(t, actualMetric.P95, time.Second)
	assert.Equal(t, actualMetric.P99, time.Second)
}


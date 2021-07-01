package metric

import "time"

func calculatePercentile(durations []time.Duration, percentile float64) time.Duration {
	return durations[int(float64(len(durations))*percentile+0.5)-1]
}

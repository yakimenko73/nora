package metric

import (
	"math"
	"time"
)

func calculatePercentile(durations []time.Duration, percentile float64) time.Duration {
	x := percentile/100 * float64(len(durations) - 1) + 1

	_, f := math.Modf(x)

	cur := float64(durations[int(x)-1])
	next := float64(durations[int(x)])
	return time.Duration(cur + f*(next - cur))
}

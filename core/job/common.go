package job

import (
	"github.com/illatior/load-testing/core/metric"
	"time"
)

type Job interface {
	Complete(start time.Time) *metric.Result
}

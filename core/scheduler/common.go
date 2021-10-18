package scheduler

import (
	"time"
)

type Scheduler interface {
	GetNextExecution(elapsed time.Duration, hits uint64) (next time.Duration, stop bool)
}

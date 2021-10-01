package scheduler

import (
	"math"
	"time"
)

type ConstantScheduler struct {
	Frequency uint64
	Period    time.Duration
}

func (cs ConstantScheduler) GetNextExecution(elapsed time.Duration, hits uint64) (next time.Duration, stop bool) {
	if cs.Period == 0 || cs.Frequency == 0 {
		return 0, false
	}
	if cs.Period < 0 {
		return 0, true
	}

	expectedHits := cs.Frequency * uint64(elapsed/cs.Period)
	if hits < expectedHits {
		return 0, false
	}
	interval := uint64(cs.Period.Nanoseconds()) / cs.Frequency
	if math.MaxInt64/interval < hits {
		return 0, true
	}

	nextElapsed := time.Duration(interval * (hits + 1))
	return nextElapsed - elapsed, false
}

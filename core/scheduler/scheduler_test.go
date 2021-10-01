package scheduler

import (
	"testing"
	"time"
)

func TestConstantSchedulerNext(t *testing.T) {
	t.Parallel()

	for ti, td := range []struct {
		freq    uint64
		period  time.Duration
		elapsed time.Duration
		hits    uint64

		next time.Duration
		stop bool
	}{
		{1, time.Second, 0, 0, time.Second, false},
		{1, time.Second, time.Second, 1, time.Second, false},
		{1, time.Second, time.Second, 2, 2 * time.Second, false},
		{1, time.Second, time.Second, 0, 0, false},

		// Infinite
		{1, 0, time.Second, 2, 0, false},
		{0, time.Second, time.Second, 2, 0, false},

		// Negative cases
		{1, -1, time.Second, 0, 0, true},
	} {
		cs := ConstantScheduler{
			Frequency: td.freq,
			Period:    td.period,
		}

		next, stop := cs.GetNextExecution(td.elapsed, td.hits)
		if next != td.next || stop != td.stop {
			t.Errorf("Test %d: (%d, %s)\nExpected: (%s, %t)\nActual: (%s, %t)",
				ti, td.freq, td.period, td.next, td.stop, next, stop)
		}
	}
}

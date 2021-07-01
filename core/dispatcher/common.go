package dispatcher

import (
	"context"
	"github.com/illatior/load-testing/core/executor"
	"github.com/illatior/load-testing/core/metric"
	"time"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, scheduler Scheduler, executor executor.Executor, duration time.Duration) <-chan *metric.Result
}

type Scheduler interface {
	ScheduleNextExecution(elapsed time.Duration, hits uint64) (next time.Duration, stop bool)
}

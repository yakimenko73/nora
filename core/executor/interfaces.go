package executor

import (
	"context"
	"github.com/illatior/nora/core/metric"
	"github.com/illatior/nora/core/task"
)

type Executor interface {
	AddTask(task task.Task)
	ScheduleExecution(ctx context.Context, ticks <-chan interface{}, results chan<- *metric.Result)
}

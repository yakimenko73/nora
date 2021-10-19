package executor

import (
	"context"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/task"
)

type Executor interface {
	AddTask(task task.Task)
	ScheduleExecution(ctx context.Context, ticks <-chan interface{}, results chan<- *metric.Result)
}

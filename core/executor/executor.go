package executor

import (
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/task"
	"sync"
)

type Executor interface {
	AddTask(task *task.Task)
	ScheduleExecution(wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result)
}

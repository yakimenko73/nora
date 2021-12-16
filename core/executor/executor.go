package executor

import (
	"github.com/yakimenko73/nora/core/metric"
	"github.com/yakimenko73/nora/core/task"
	"sync"
)

type Executor interface {
	AddTask(task *task.Task)
	ScheduleExecution(wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result)
}

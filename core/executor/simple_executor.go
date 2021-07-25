package executor

import (
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/task"
	"sync"
)

type simpleExecutor struct {
	tasks []*task.Task
}

func New() Executor {
	return &simpleExecutor{
		tasks: make([]*task.Task, 0),
	}
}

func (e *simpleExecutor) AddTask(task *task.Task) {
	e.tasks = append(e.tasks, task)
}

func (e *simpleExecutor) ScheduleExecution(wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result) {
	defer wg.Done()

	for range ticks {
		for _, j := range e.tasks {
			results <- (*j).Run()
		}
	}
}

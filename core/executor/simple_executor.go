package executor

import (
	"context"
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

func (e *simpleExecutor) ScheduleExecution(ctx context.Context, ticks <-chan interface{}, results chan<- *metric.Result) {
	var wg sync.WaitGroup
	childCtx, cancel := context.WithCancel(ctx)

	defer wg.Wait()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			for _, t := range e.tasks {

				wg.Add(1)
				go func(t *task.Task) {
					defer wg.Done()
					results <- (*t).Run(childCtx)
				}(t)
			}
		}
	}
}

package executor

import (
	"context"
	"errors"
	"github.com/illatior/nora/lib/metric"
	"github.com/illatior/nora/lib/task"
	"github.com/illatior/nora/lib/util"
	"sync"
)

type randomExecutor struct {
	tasks []task.Task
}

func NewRandomExecutor() Executor {
	return &randomExecutor{
		tasks: make([]task.Task, 0),
	}
}

func (re *randomExecutor) AddTask(task task.Task) {
	re.tasks = append(re.tasks, task)
}

// ScheduleExecution method is blocking
func (re *randomExecutor) ScheduleExecution(ctx context.Context, ticks <-chan interface{}, results chan<- *metric.Result) {
	if len(re.tasks) == 0 {
		panic(errors.New("executor not initialized"))
	}

	var wg sync.WaitGroup

	childCtx, cancel := context.WithCancel(ctx)
	defer wg.Wait()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			nextJobIndex, err := util.GetRandomInt(0, len(re.tasks))
			if err != nil {
				panic(err)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- re.tasks[nextJobIndex].Run(childCtx)
			}()
		}
	}
}

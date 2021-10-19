package executor

import (
	"context"
	"errors"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/task"
	"sync"
	"sync/atomic"
)

type roundRobinExecutor struct {
	tasks []*task.Task

	next uint64
}

func NewRoundRobinExecutor() Executor {
	return &roundRobinExecutor{
		tasks: make([]*task.Task, 0),
		next:  0,
	}
}

func (e *roundRobinExecutor) AddTask(task *task.Task) {
	e.tasks = append(e.tasks, task)
}

// ScheduleExecution method is blocking
func (e *roundRobinExecutor) ScheduleExecution(ctx context.Context, ticks <-chan interface{}, results chan<- *metric.Result) {
	var wg sync.WaitGroup

	childCtx, cancel := context.WithCancel(ctx)
	defer wg.Wait()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			next, err := e.getNext()
			if err != nil {
				panic(err)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				results <- (*e.tasks[next]).Run(childCtx)
			}()
		}
	}
}

func (e *roundRobinExecutor) getNext() (uint64, error) {
	if len(e.tasks) == 0 {
		return 0, errors.New("executor not initialized (tasks == 0)")
	}

	atomic.AddUint64(&e.next, 1)
	return e.next % uint64(len(e.tasks)), nil
}

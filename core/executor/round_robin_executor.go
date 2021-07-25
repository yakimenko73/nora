package executor

import (
	"errors"
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/task"
	"sync"
	"sync/atomic"
)

type roundRobinExecutor struct {
	tasks []*task.Task

	next int64
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

func (e *roundRobinExecutor) ScheduleExecution(wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result) {
	defer wg.Done()

	for range ticks {
		next, err := e.getNext()
		if err != nil {
			panic(err)
		}

		results <- (*e.tasks[next]).Run()
	}
}

func (e *roundRobinExecutor) getNext() (int64, error) {
	if len(e.tasks) == 0 {
		return 0, errors.New("Executor not initialized (tasks == 0)!")
	}

	atomic.AddInt64(&e.next, 1)
	if e.next == int64(len(e.tasks)) {
		e.next = 0
	}

	return e.next, nil
}

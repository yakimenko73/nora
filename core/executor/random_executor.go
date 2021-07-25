package executor

import (
	"github.com/illatior/task-scheduler/core/metric"
	"github.com/illatior/task-scheduler/core/task"
	"github.com/illatior/task-scheduler/core/util"
	"sync"
)

type randomExecutor struct {
	tasks []*task.Task
}

func NewRandomExecutor() Executor {
	return &randomExecutor{
		tasks: make([]*task.Task, 0),
	}
}

func (re *randomExecutor) AddTask(task *task.Task) {
	re.tasks = append(re.tasks, task)
}

func (re *randomExecutor) ScheduleExecution(wg *sync.WaitGroup, ticks <-chan interface{}, results chan<- *metric.Result) {
	defer wg.Done()

	for range ticks {
		nextJobIndex, err := util.GetRandomInt(0, len(re.tasks))
		if err != nil {
			panic(err) // FIXME dont panic here
		}

		results <- (*re.tasks[nextJobIndex]).Run()
	}
}

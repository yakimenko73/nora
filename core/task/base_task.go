package task

import (
	"github.com/yakimenko73/nora/core/metric"
	"time"
)

type baseTask struct {
	task func() error
	name string
}

func NewBaseTask(taskFunc func() error, name string) Task {
	return &baseTask{
		task: taskFunc,
		name: name,
	}
}

func (t *baseTask) Run() (res *metric.Result) {
	var err error
	res = &metric.Result{
		Name:  t.name,
		Start: time.Now(),
	}
	defer func() {
		if err != nil {
			res.Error = err
		}
		res.Duration = time.Since(res.Start)
	}()

	err = t.task()
	return
}

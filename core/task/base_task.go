package task

import (
	"context"
	"github.com/illatior/nora/core/metric"
	"time"
)

type baseTask struct {
	task func(ctx context.Context) error
	name string
}

func NewBaseTask(taskFunc func(ctx context.Context) error, name string) Task {
	return &baseTask{
		task: taskFunc,
		name: name,
	}
}

func (t *baseTask) Run(ctx context.Context) (res *metric.Result) {
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

	err = t.task(ctx)
	return
}

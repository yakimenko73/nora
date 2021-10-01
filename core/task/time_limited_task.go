package task

import (
	"context"
	"github.com/illatior/task-scheduler/core/metric"
	"time"
)

type timeLimitedTask struct {
	task func(ctx context.Context) error
	name string

	timeout time.Duration
}

func NewTimeLimitedTask(taskFunc func(ctx context.Context) error, name string, timeout time.Duration) Task {
	return &timeLimitedTask{
		task:    taskFunc,
		name:    name,
		timeout: timeout,
	}
}

func (t *timeLimitedTask) Run(ctx context.Context) (res *metric.Result) {
	var err error
	childCtx, cancel := context.WithTimeout(ctx, t.timeout)

	done := make(chan bool, 1)

	go func() {
		defer cancel()

		select {
		case <-ctx.Done():
			err = errContextCancelled
			return
		case <-done:
			return
		case <-time.After(t.timeout):
			err = errTaskTimedOut
			return
		}
	}()

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

	err = t.task(childCtx)
	done <- true
	return
}

package task

import (
	"context"
	"github.com/illatior/nora/lib/metric"
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
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

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

	resCh := make(chan error, 1)
	go func() {
		resCh <- t.task(ctx)
	}()

	select {
	case <-ctx.Done():
		err = errContextCancelled
	case <-time.After(t.timeout):
		err = errTaskTimedOut
	case res := <-resCh:
		err = res
	}

	return
}

package task

import (
	"context"
	"github.com/illatior/nora/core/metric"
	"github.com/pkg/errors"
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
	childCtx, cancel := context.WithTimeout(ctx, t.timeout)

	done := make(chan bool, 1)

	var err error
	go func() {
		defer cancel()

		select {
		case <-done:
			return
		case <-ctx.Done():
			err = errContextCancelled
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

	if tmp := t.task(childCtx); tmp != nil {
		if err != nil {
			err = errors.Wrap(err, tmp.Error())
		} else {
			err = tmp
		}
	}
	done <- true
	return
}

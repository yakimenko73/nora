package task

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBaseTask_Run(t *testing.T) {
	t.Parallel()

	isExecuted := false

	taskName := "test"
	taskFunc := func(ctx context.Context) error {
		isExecuted = true
		return nil
	}
	start := time.Now()
	delta := time.Nanosecond

	task := NewBaseTask(taskFunc, taskName)

	res := task.Run(context.Background())

	assert.True(t, isExecuted)
	assert.Equal(t, res.Name, taskName)
	assert.Nil(t, res.Error)

	assert.GreaterOrEqual(t, res.Duration, time.Duration(0))
	assert.LessOrEqual(t, res.Duration, res.Start.Add(delta).Sub(res.Start))

	assert.GreaterOrEqual(t, res.Start.Sub(start), time.Duration(0))
	assert.LessOrEqual(t, res.Start.Sub(start), delta)
}

func TestBaseTask_ContextCancelled(t *testing.T) {
	t.Parallel()

	start := time.Now()

	delta := 100 * time.Millisecond
	contextTimeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	isHandledContextDone := false

	taskName := "test"
	taskFunc := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			isHandledContextDone = true
			return nil
		case <-time.After(contextTimeout * 2):
			return nil
		}
	}

	taskk := NewBaseTask(taskFunc, taskName)
	taskk.Run(ctx)

	assert.True(t, isHandledContextDone)

	since := time.Since(start)
	assert.GreaterOrEqual(t, since, time.Duration(0))
	assert.LessOrEqual(t, since, contextTimeout+delta)
}

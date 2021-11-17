package task

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeLimitedTask_Run(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	start := time.Now()

	maxExecutionTime := 5 * time.Microsecond

	taskName := "test task"
	taskTimeout := 1 * time.Second

	isExecuted := false
	taskFunction := func(ctx context.Context) error {
		isExecuted = true
		return nil
	}

	taskk := NewTimeLimitedTask(taskFunction, taskName, taskTimeout)
	res := taskk.Run(ctx)

	assert.True(t, isExecuted)

	since := time.Since(start)
	assert.GreaterOrEqual(t, since, time.Duration(0))
	assert.LessOrEqual(t, since, maxExecutionTime)

	assert.Equal(t, res.Name, taskName)
	assert.Nil(t, res.Error)

	assert.GreaterOrEqual(t, res.Start.Sub(start), time.Duration(0))
	assert.LessOrEqual(t, res.Start.Sub(start), since)

	assert.GreaterOrEqual(t, res.Duration, time.Duration(0))
	assert.LessOrEqual(t, res.Duration, since)
}

func TestTimeLimitedTask_TimeOut(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	start := time.Now()

	maxExecutionTime := 5 * time.Microsecond

	taskName := "test task"
	taskTimeout := 1 * time.Second

	isExecuted := false
	taskFunction := func(ctx context.Context) error {
		isExecuted = true
		time.Sleep(taskTimeout + 50*time.Millisecond)
		return nil
	}

	taskk := NewTimeLimitedTask(taskFunction, taskName, taskTimeout)
	res := taskk.Run(ctx)

	assert.True(t, isExecuted)

	since := time.Since(start)
	assert.GreaterOrEqual(t, since, time.Duration(0))
	assert.LessOrEqual(t, since, maxExecutionTime)

	assert.Equal(t, res.Name, taskName)
	assert.Equal(t, res.Error, errTaskTimedOut)

	assert.GreaterOrEqual(t, res.Start.Sub(start), time.Duration(0))
	assert.LessOrEqual(t, res.Start.Sub(start), since)

	assert.GreaterOrEqual(t, res.Duration, time.Duration(0))
	assert.LessOrEqual(t, res.Duration, since)
}

func TestTimeLimitedTask_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	taskFunction := func(ctx context.Context) error {
		cancel()
		time.Sleep(10 * time.Millisecond)

		return nil
	}

	taskk := NewTimeLimitedTask(taskFunction, "", time.Minute)
	res := taskk.Run(ctx)

	assert.Equal(t, res.Error, errContextCancelled)
}

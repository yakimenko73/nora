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

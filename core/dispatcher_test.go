package core

import (
	"context"
	"fmt"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/scheduler"
	"github.com/illatior/task-scheduler/core/task"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestCallDispatch(t *testing.T) {
	t.Parallel()
	testDuration := time.Second * 2

	deltaBelow := 3
	deltaAbove := 1
	dispatchDuration := 1 * time.Second

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(testDuration))
	defer cancel()

	sch := scheduler.ConstantScheduler{
		Frequency: 1,
		Period:    1 * time.Microsecond,
	}
	taskk := task.NewBaseTask(func(ctx context.Context) error {
		return nil
	}, "test")

	exec := executor.New()
	exec.AddTask(&taskk)

	actualExecutions := uint64(0)
	for range Dispatch(ctx, sch, exec, dispatchDuration, uint64(runtime.GOMAXPROCS(0))) {
		actualExecutions++
	}
	t.Log(fmt.Sprintf("Total executions: %d", actualExecutions))

	if ctx.Err() != nil {
		t.Fatal(ctx.Err())
	}

	expectedTotalExecutions := uint64(dispatchDuration/sch.Period) * sch.Frequency

	expectedMinExecutions := expectedTotalExecutions * uint64(100-deltaBelow) / 100
	expectedMaxExecutions := expectedTotalExecutions * uint64(100+deltaAbove) / 100
	assert.GreaterOrEqual(t, actualExecutions, expectedMinExecutions)
	assert.LessOrEqual(t, actualExecutions, expectedMaxExecutions)
}

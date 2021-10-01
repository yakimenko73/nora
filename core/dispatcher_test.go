package core

import (
	"context"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/scheduler"
	"github.com/illatior/task-scheduler/core/task"
	"runtime"
	"testing"
	"time"
)

func TestCallDispatch(t *testing.T) {
	t.Parallel()

	sch := scheduler.ConstantScheduler{
		Frequency: 1,
		Period:    1 * time.Millisecond,
	}

	taskk := task.NewBaseTask(func() error {
		return nil
	}, "task name")

	exec := executor.New()
	exec.AddTask(&taskk)

	for range Dispatch(context.Background(), sch, exec, time.Second, uint64(runtime.GOMAXPROCS(0))) {
	}
}

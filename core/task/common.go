package task

import (
	"context"
	"errors"
	"github.com/illatior/task-scheduler/core/metric"
)

type Task interface {
	Run(ctx context.Context) *metric.Result
}

var (
	errTaskTimedOut     = errors.New("task timed out")
	errContextCancelled = errors.New("parent context cancelled")
)

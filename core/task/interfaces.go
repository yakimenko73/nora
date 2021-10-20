package task

import (
	"context"
	"github.com/illatior/task-scheduler/core/metric"
)

type Task interface {
	Run(ctx context.Context) *metric.Result
}

package task

import (
	"context"
	"github.com/illatior/nora/core/metric"
)

type Task interface {
	Run(ctx context.Context) *metric.Result
}

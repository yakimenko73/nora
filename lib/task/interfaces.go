package task

import (
	"context"
	"github.com/illatior/nora/lib/metric"
)

type Task interface {
	Run(ctx context.Context) *metric.Result
}

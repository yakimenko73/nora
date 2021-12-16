package task

import (
	"github.com/yakimenko73/nora/core/metric"
)

type Task interface {
	Run() *metric.Result
}

package task

import (
	"github.com/illatior/task-scheduler/core/metric"
)

type Task interface {
	Run() *metric.Result
}

package job

import (
	"github.com/illatior/load-testing/core/metric"
	"time"
)

type baseJob struct {
	job  func() error
	name string
}

func NewBaseJob(job func() error, name string) Job {
	return &baseJob{
		job:  job,
		name: name,
	}
}

func (j *baseJob) Complete(start time.Time) *metric.Result {
	var err error
	res := &metric.Result{
		Name:  j.name,
		Start: start,
	}
	defer func() {
		if err != nil {
			res.Error = err
		}
		res.Duration = time.Since(start)
	}()

	err = j.job()

	return res
}

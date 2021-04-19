package worker

import "load-testing/core/job"

type Worker interface {
	Run() error
	Stop()
	AddJob(job job.Job) error
}

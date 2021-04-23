package load

import (
	"load-testing/core/job"
	"load-testing/core/worker"
)

type LoadService interface {
	Start()
	AddJob(jobFunc func() error) error
	AddJobToSpecificWorker(jobFunc func() error, workerType worker.WorkerType, jobType job.JobType, appendToPrevious bool) error
	SetLoadTime(loadTime float32)
}

package job

import "errors"

type JobType uint16

const (
	BaseJob = iota << 1
	Unknown
)

func Classify(jobType JobType, job func() error) (Job, error) {
	switch jobType {
	case BaseJob:
		return NewBaseJob(job), nil
	default:
		return nil, errors.New("Unable to classify worker type")
	}
}

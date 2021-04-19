package job

type baseJob struct {
	job func() error
}

func NewBaseJob(job func() error) Job {
	return &baseJob{
		job: job,
	}
}

func (j *baseJob) Complete() error {
	return j.job()
}

func (j *baseJob) Cancel() {
	return
}

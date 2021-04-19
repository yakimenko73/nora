package worker

import (
	"errors"
	"load-testing/core/job"
	"sync"

	"golang.org/x/sync/errgroup"
)

type waitableWorker struct {
	stopped bool

	mtx sync.RWMutex

	jobs []job.Job
}

func NewWaitableWorker() Worker {
	return &waitableWorker{
		stopped: false,
		mtx:     sync.RWMutex{},
	}
}

func (w *waitableWorker) Run() error {
	if w.stopped {
		return errors.New("Worker is currently stopped")
	}
	var eg errgroup.Group

	for _, job := range w.jobs {
		job := job
		eg.Go(func() error {
			return job.Complete()
		})
	}

	return eg.Wait()
}

func (w *waitableWorker) Stop() {
	w.stopped = true
}

func (w *waitableWorker) AddJob(job job.Job) error {
	if w.stopped {
		return errors.New("Worker already stopped")
	}
	w.mtx.Lock()
	defer w.mtx.Unlock()

	w.jobs = append(w.jobs, job)

	return nil
}

package load

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"load-testing/core/dispatcher"
	"load-testing/core/job"
	"load-testing/core/metric"
	"load-testing/core/worker"
	"time"
)

type baseLoadService struct {
	loadTime float32

	dispatcher        dispatcher.Dispatcher
	currentWorker     worker.Worker
	currentWorkerType worker.WorkerType

	ctx            context.Context
	metricConsumer *metric.MetricConsumer
}

func NewLoadService(dispatcher dispatcher.Dispatcher, ctx context.Context) LoadService {
	return &baseLoadService{
		loadTime:   0,
		dispatcher: dispatcher,
		ctx:        ctx,
	}
}

func (ls *baseLoadService) Start() {
	ctx, cancel := context.WithCancel(ls.ctx)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ls.dispatcher.Dispatch(ctx, ls.metricConsumer)
	})

	g.Go(func() error {
		select {
		case <-time.After(time.Duration(float32(time.Minute) * ls.loadTime)):
			cancel()
			return nil
		}
	})

	g.Wait()
}

// TODO
func (ls *baseLoadService) AddJobToSpecificWorker(jobFunc func() error, workerType worker.WorkerType, jobType job.JobType, appendToPrevious bool) error {
	jobObj, err := job.Classify(jobType, jobFunc)
	if err != nil {
		return err
	}

	if ls.currentWorker == nil {
		worker, err := worker.Classify(workerType)
		if err != nil {
			return err
		}

		ls.currentWorker = worker
		ls.currentWorkerType = workerType
		err = ls.dispatcher.AddWorker(uuid.NewV4().String(), &ls.currentWorker)
		if err != nil {
			return err
		}
	}

	if workerType == ls.currentWorkerType && appendToPrevious {
		ls.currentWorker.AddJob(jobObj)
	} else {

		worker, err := worker.Classify(workerType)
		if err != nil {
			return err
		}

		ls.currentWorker = worker
		ls.currentWorkerType = workerType
		ls.currentWorker.AddJob(jobObj)

		err = ls.dispatcher.AddWorker(uuid.NewV4().String(), &ls.currentWorker)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ls *baseLoadService) AddJob(jobFunc func() error) error {
	return ls.AddJobToSpecificWorker(jobFunc, ls.currentWorkerType, job.BaseJob, true)
}

func (ls *baseLoadService) SetLoadTime(loadTime float32) {
	ls.loadTime = loadTime
}

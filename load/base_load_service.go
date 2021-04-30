package load

import (
	"context"
	"golang.org/x/sync/errgroup"
	"load-testing/core/dispatcher"
	"load-testing/core/job"
	"load-testing/core/metric"
	"math/rand"
	"strconv"
	"time"
)

type baseLoadService struct {
	loadTime float32

	dispatcher        dispatcher.Dispatcher

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
	ctx, _ := context.WithTimeout(ls.ctx, time.Duration(float32(time.Minute) * ls.loadTime))
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ls.dispatcher.Dispatch(ctx, ls.metricConsumer)
	})

	_ = g.Wait()
}

func (ls *baseLoadService) AddJob(jobFunc func() error) error {
	j := job.NewBaseJob(jobFunc)
	return ls.dispatcher.AddJob(strconv.Itoa(rand.Int()), j)
}

func (ls *baseLoadService) SetLoadTime(loadTime float32) {
	ls.loadTime = loadTime
}

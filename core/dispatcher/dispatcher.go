package dispatcher

import (
	"context"
	"load-testing/core/worker"
)

type Dispatcher interface {
	AddWorker(id string, worker *worker.Worker) error
	Dispatch(ctx context.Context) error
	Shutdown() error
}

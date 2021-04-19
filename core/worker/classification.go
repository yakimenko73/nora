package worker

import "errors"

type WorkerType uint16

const (
	WaitableWorker = iota << 1
	Unknown
)

func Classify(workerType WorkerType) (Worker, error) {
	switch workerType {
	case WaitableWorker:
		return NewWaitableWorker(), nil
	default:
		return nil, errors.New("Unable to classify worker type")
	}
}

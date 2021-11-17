package task

import "errors"

var (
	errTaskTimedOut     = errors.New("task timed out")
	errContextCancelled = errors.New("parent context cancelled")
)

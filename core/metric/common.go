package metric

import "time"

type Result struct {
	Name string

	Start    time.Time
	Duration time.Duration

	Error           error
}

package core

import "time"

type LoadOptions struct {
	Duration time.Duration
	Workers  int

	Frequency uint64 // ticks per ....
	Period    time.Duration
}

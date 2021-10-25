package core

import (
	"context"
	"errors"
	"github.com/illatior/task-scheduler/core/executor"
	"github.com/illatior/task-scheduler/core/task"
	"time"
)

type Option interface {
	apply(d *Dispatcher) error
}

type option func(d *Dispatcher) error

func (o option) apply(d *Dispatcher) error {
	return o(d)
}

func WithDuration(dur time.Duration) Option {
	return option(func(d *Dispatcher) error {
		if dur < 0 {
			return errors.New("load duration should be positive")
		}

		d.configuration.Duration = dur
		return nil
	})
}

func WithPeriod(p time.Duration) Option {
	return option(func(d *Dispatcher) error {
		if p < 0 {
			return errors.New("period should be positive")
		}

		d.configuration.Period = p
		return nil
	})
}

func WithFrequency(f uint64) Option {
	return option(func(d *Dispatcher) error {
		if f < 0 {
			return errors.New("frequency should be positive")
		}

		d.configuration.Frequency = f
		return nil
	})
}

func WithExecutor(e executor.Executor) Option {
	return option(func(d *Dispatcher) error {

		d.exec = e
		return nil
	})
}

func WithWorkersCount(c int) Option {
	return option(func(d *Dispatcher) error {
		if c <= 0 {
			return errors.New("workers count should be greater than zero")
		}

		d.configuration.Workers = c
		return nil
	})
}

func WithTask(name string, f func(context.Context) error) Option {
	return option(func(d *Dispatcher) error {
		if name == "" {
			return errors.New("task name must not be empty")
		}

		t := task.NewBaseTask(f, name)
		d.exec.AddTask(t)
		return nil
	})
}

func WithTimeLimitedTask(name string, timeout time.Duration, f func(context.Context) error) Option {
	return option(func(d *Dispatcher) error {
		if name == "" {
			return errors.New("task name must not be empty")
		}
		if timeout <= 0 {
			return errors.New("time-limited task timeout must be > 0")
		}

		t := task.NewTimeLimitedTask(f, name, timeout)
		d.exec.AddTask(t)
		return nil
	})
}

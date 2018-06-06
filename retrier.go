package amqpx

import "time"

type retryStrategy int

const (
	retryStrategyNoop retryStrategy = iota
	retryStrategyExponential
)

type retrier interface {
	retry(func() error) error
}
type retrierOptions struct {
	strategy    retryStrategy
	exponential struct {
		initialInterval time.Duration
		maxInterval     time.Duration
		maxElapsedTime  time.Duration
	}
}

type retriersOptions struct {
	connection retrierOptions
	channel    retrierOptions
}

func newRetrier(opts retrierOptions) retrier {
	switch opts.strategy {
	case retryStrategyExponential:
		return newExponentialRetrier(opts)
	case retryStrategyNoop:
		return newNoopRetrier()
	default:
		return newNoopRetrier()
	}
}

type retriers struct {
	connection retrier
	channel    retrier
}

func newRetriers(opts retriersOptions) *retriers {
	return &retriers{
		connection: newRetrier(opts.connection),
		channel:    newRetrier(opts.channel),
	}
}

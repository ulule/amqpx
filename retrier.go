package amqpx

import "time"

type retryStrategy uint8

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

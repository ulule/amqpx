package amqpx

type retryStrategy int

const (
	retryStrategyNoop retryStrategy = iota
	retryStrategyExponential
)

type retrier interface {
	retry(func() error) error
}

func newRetrier(opts retryOptions) retrier {
	switch opts.strategy {
	case retryStrategyExponential:
		return newExponentialRetrier(opts)
	case retryStrategyNoop:
		return newNoopRetrier()
	default:
		return newNoopRetrier()
	}
}

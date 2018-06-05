package amqpx

type retryStrategy int

const (
	retryStrategyNoop retryStrategy = iota
	retryStrategyExponential
)

type retrier interface {
	retry(func() error) error
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

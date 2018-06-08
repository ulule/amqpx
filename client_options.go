package amqpx

import "time"

// ClientOption is used to define Client options.
type ClientOption interface {
	apply(*clientOptions) error
}

type clientOption func(*clientOptions) error

func (o clientOption) apply(instance *clientOptions) error {
	return o(instance)
}

type clientOptions struct {
	dialer   Dialer
	observer Observer
	logger   Logger
	usePool  bool
	capacity int
	retrier  retrierOptions
}

// WithCapacity will configure a Client with the given number of connections.
func WithCapacity(capacity int) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if capacity <= 0 {
			return ErrInvalidConnectionsPoolCapacity
		}
		options.usePool = true
		options.capacity = capacity
		return nil
	})
}

// WithoutConnectionsPool will configure a Client without a connections pool.
func WithoutConnectionsPool() ClientOption {
	return clientOption(func(options *clientOptions) error {
		options.usePool = false
		return nil
	})
}

// WithObserver will configure Client with the given observer.
func WithObserver(observer Observer) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if observer == nil {
			return ErrObserverRequired
		}
		options.observer = observer
		return nil
	})
}

// WithLogger will configure Client with the given Logger.
func WithLogger(logger Logger) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if logger == nil {
			return ErrLoggerRequired
		}
		options.logger = logger
		return nil
	})
}

// WithDefaultLogger will configure Client with the defaut logger.
func WithDefaultLogger(level LoggerLevel) ClientOption {
	return clientOption(func(options *clientOptions) error {
		options.logger = newDefaultLogger(level)
		return nil
	})
}

// WithExponentialRetry will configure a Client with the given exponential
// connection retry durations.
func WithExponentialRetry(initialInterval, maxInterval, maxElapsedTime time.Duration) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if initialInterval <= 0 || maxInterval <= 0 || maxElapsedTime <= 0 {
			return ErrInvalidRetryDuration
		}
		options.retrier.strategy = retryStrategyExponential
		options.retrier.exponential.initialInterval = initialInterval
		options.retrier.exponential.maxInterval = maxInterval
		options.retrier.exponential.maxElapsedTime = maxElapsedTime
		return nil
	})
}

package amqpx

import "time"

// ClientOption is used to define client configuration.
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
	usePool  bool
	capacity int
	retriers retriersOptions
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

// WithCapacity will configure a client with a connections pool and given capacity.
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

// WithoutConnectionsPool will configure a client without a connections pool.
func WithoutConnectionsPool() ClientOption {
	return clientOption(func(options *clientOptions) error {
		options.usePool = false
		return nil
	})
}

// WithObserver will configure client with given observer.
func WithObserver(observer Observer) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if observer == nil {
			return ErrObserverRequired
		}

		options.observer = observer
		return nil
	})
}

// WithExponentialConnectionRetry will configure a client with the given exponential
// connection retry durations.
func WithExponentialConnectionRetry(initialInterval, maxInterval, maxElapsedTime time.Duration) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if initialInterval <= 0 || maxInterval <= 0 || maxElapsedTime <= 0 {
			return ErrInvalidRetryDuration
		}

		options.retriers.connection.strategy = retryStrategyExponential
		options.retriers.connection.exponential.initialInterval = initialInterval
		options.retriers.connection.exponential.maxInterval = maxInterval
		options.retriers.connection.exponential.maxElapsedTime = maxElapsedTime

		return nil
	})
}

// WithExponentialChannelRetry will configure a client with the given exponential
// channel retry durations.
func WithExponentialChannelRetry(initialInterval, maxInterval, maxElapsedTime time.Duration) ClientOption {
	return clientOption(func(options *clientOptions) error {
		if initialInterval <= 0 || maxInterval <= 0 || maxElapsedTime <= 0 {
			return ErrInvalidRetryDuration
		}

		options.retriers.channel.strategy = retryStrategyExponential
		options.retriers.channel.exponential.initialInterval = initialInterval
		options.retriers.channel.exponential.maxInterval = maxInterval
		options.retriers.channel.exponential.maxElapsedTime = maxElapsedTime

		return nil
	})
}

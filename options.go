package amqpx

import "time"

// Option is used to define client configuration.
type Option interface {
	apply(*clientOptions) error
}

type option func(*clientOptions) error

func (o option) apply(instance *clientOptions) error {
	return o(instance)
}

type clientOptions struct {
	dialer   Dialer
	observer Observer
	usePool  bool
	capacity int
	retryOptions
}

type retryOptions struct {
	useRetry             bool
	retryInitialInterval time.Duration
	retryMaxInterval     time.Duration
	retryMaxElapsedTime  time.Duration
}

// WithCapacity will configure a client with a connections pool and given capacity.
func WithCapacity(capacity int) Option {
	return option(func(options *clientOptions) error {
		if capacity <= 0 {
			return ErrInvalidConnectionsPoolCapacity
		}

		options.usePool = true
		options.capacity = capacity
		return nil
	})
}

// WithoutConnectionsPool will configure a client without a connections pool.
func WithoutConnectionsPool() Option {
	return option(func(options *clientOptions) error {
		options.usePool = false
		return nil
	})
}

// WithObserver will configure client with given observer.
func WithObserver(observer Observer) Option {
	return option(func(options *clientOptions) error {
		if observer == nil {
			return ErrObserverRequired
		}

		options.observer = observer
		return nil
	})
}

// WithRetry will enable retry.
func WithRetry() Option {
	return option(func(options *clientOptions) error {
		options.useRetry = true
		return nil
	})
}

// WithRetryDurations will configure a client with the given retry durations.
func WithRetryDurations(initialInterval, maxInterval, maxElapsedTime time.Duration) Option {
	return option(func(options *clientOptions) error {
		if initialInterval <= 0 || maxInterval <= 0 || maxElapsedTime <= 0 {
			return ErrInvalidRetryDuration
		}

		options.useRetry = true
		options.retryInitialInterval = initialInterval
		options.retryMaxInterval = maxInterval
		options.retryMaxElapsedTime = maxElapsedTime

		return nil
	})
}

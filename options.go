package amqpx

// Option is used to define client configuration.
type Option interface {
	apply(*clientOptions) error
}

type option func(*clientOptions) error

func (o option) apply(instance *clientOptions) error {
	return o(instance)
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

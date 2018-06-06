package amqpx

import "time"

// DialerOption is used to define dialer configuration.
type DialerOption interface {
	apply(*dialerOptions) error
}

type dialerOption func(*dialerOptions) error

func (o dialerOption) apply(instance *dialerOptions) error {
	return o(instance)
}

type dialerOptions struct {
	timeout   time.Duration
	heartbeat time.Duration
}

func newDialerOptions() dialerOptions {
	return dialerOptions{
		timeout:   DefaultDialerTimeout,
		heartbeat: DefaultDialerHeartbeat,
	}
}

// WithDialerTimeout will configure a dialer with the given timeout duration.
func WithDialerTimeout(timeout time.Duration) DialerOption {
	return dialerOption(func(options *dialerOptions) error {
		if timeout <= 0 {
			return ErrInvalidDialerTimeout
		}
		options.timeout = timeout
		return nil
	})
}

// WithDialerHeartbeat will configure a dialer with the given heartbeat duration.
func WithDialerHeartbeat(heartbeat time.Duration) DialerOption {
	return dialerOption(func(options *dialerOptions) error {
		if heartbeat <= 0 {
			return ErrInvalidDialerHeartbeat
		}
		options.heartbeat = heartbeat
		return nil
	})
}

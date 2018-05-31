// Package amqpx provides extensions for Go AMQP library.
//
// It use https://github.com/streadway/amqp as dependency.
package amqpx

import (
	"time"

	"github.com/pkg/errors"
)

const (
	defaultRetryInitialInterval = 100 * time.Millisecond
	defaultRetryMaxInterval     = 32 * time.Second
	defaultRetryMaxElapsedTime  = 7 * time.Minute
)

// Client interface describe a amqp client.
type Client interface {
	// Channel returns a new amqp's channel from current client unless it's closed.
	Channel() (Channel, error)

	// Close closes the client.
	Close() error

	// IsClosed returns if the client is closed.
	IsClosed() bool
}

// New returns a new client using given configuration.
func New(dialer Dialer, options ...Option) (Client, error) {
	// Default client options.
	opts := &clientOptions{
		dialer:   dialer,
		observer: &defaultObserver{},
		usePool:  true,
		capacity: DefaultConnectionsCapacity,

		// By default, retry is disabled but we set sane defaults.
		retryOptions: retryOptions{
			useRetry:             false,
			retryInitialInterval: defaultRetryInitialInterval,
			retryMaxInterval:     defaultRetryMaxInterval,
			retryMaxElapsedTime:  defaultRetryMaxElapsedTime,
		},
	}

	// Applies options.
	for _, option := range options {
		err := option.apply(opts)
		if err != nil {
			return nil, errors.Wrap(err, "amqpx: cannot create a new client")
		}
	}

	if !opts.usePool {
		return newSimple(opts)
	}

	return newConnectionsPool(opts)
}

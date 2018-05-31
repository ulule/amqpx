// Package amqpx provides extensions for Go AMQP library.
//
// It use https://github.com/streadway/amqp as dependency.
package amqpx

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	defaultRetryInitialInterval = 100 * time.Millisecond
	defaultRetryMaxInterval     = 32 * time.Second
	defaultRetryMaxElapsedTime  = 7 * time.Minute
)

type (
	// Table is an alias to amqp.Table
	Table = amqp.Table
	// Queue is an alias to amqp.Queue
	Queue = amqp.Queue
	// Publishing is an alias to amqp.Publishing
	Publishing = amqp.Publishing
	// Delivery is an alias to amqp.Delivery
	Delivery = amqp.Delivery
	// Confirmation is an alias to amqp.Confirmation
	Confirmation = amqp.Confirmation
	// Return is an alias to amqp.Return
	Return = amqp.Return
	// Error is an alias to amqp.Error
	Error = amqp.Error
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

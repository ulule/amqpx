// Package amqpx provides extensions for Go AMQP library.
//
// It uses https://github.com/streadway/amqp as dependency.
package amqpx

import (
	"github.com/pkg/errors"
)

// Client interface describes a amqp client.
type Client interface {
	// Channel returns a new amqp's channel from current client unless it's closed.
	Channel() (Channel, error)

	// Close closes the client.
	Close() error

	// IsClosed returns if the client is closed.
	IsClosed() bool
}

// New returns a new Client with the given Dialer and options.
func New(dialer Dialer, options ...ClientOption) (Client, error) {
	opts := &clientOptions{
		dialer:   dialer,
		observer: &defaultObserver{},
		usePool:  true,
		capacity: DefaultConnectionsCapacity,
		retriers: retriersOptions{},
	}

	for _, option := range options {
		err := option.apply(opts)
		if err != nil {
			return nil, errors.Wrap(err, ErrMessageCannotCreateClient)
		}
	}

	if !opts.usePool {
		return newSimple(opts)
	}

	return newPool(opts)
}

// Package amqpx provides extensions for Go AMQP library.
//
// It uses https://github.com/streadway/amqp as dependency.
package amqpx

import (
	"fmt"

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

	// channel returns a new amqp's channel from current client unless it's closed.
	channel() (*amqp.Channel, error)
}

// New returns a new Client with the given Dialer and options.
func New(dialer Dialer, options ...ClientOption) (Client, error) {
	opts := &clientOptions{
		dialer:   dialer,
		observer: &defaultObserver{},
		logger:   &noopLogger{},
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
		opts.logger.Debug("Connection pooling is disabled")
		return newSimple(opts)
	}

	opts.logger.Debug(fmt.Sprintf("Connection pooling is enabled (%d connections)", opts.capacity))

	return newPool(opts)
}

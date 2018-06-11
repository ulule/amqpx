package amqpx

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Client interface describes a amqp client.
type Client interface {
	// Channel returns a new Channel from current client unless it's closed.
	Channel() (*amqp.Channel, error)

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
		logger:   &noopLogger{},
		usePool:  true,
		capacity: DefaultConnectionsCapacity,
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

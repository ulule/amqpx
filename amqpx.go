// Package amqpx provides extensions for Go AMQP library.
//
// It use https://github.com/streadway/amqp as dependency.
package amqpx

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var (
	// ErrClientClosed occurs when operating on a closed client.
	ErrClientClosed = fmt.Errorf("client is closed")
)

// Dialer is a function returning a new amqp connection.
type Dialer func() (*amqp.Connection, error)

// Client interface describe a amqp client.
type Client interface {

	// Channel returns a new amqp's channel from current client unless it's closed.
	Channel() (*amqp.Channel, error)

	// Close closes the client.
	Close() error

	// IsClosed returns if the client is closed.
	IsClosed() bool
}

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
}

// New returns a new client using given configuration.
func New(dialer Dialer, options ...Option) (Client, error) {
	// Default client options.
	opts := &clientOptions{
		dialer:   dialer,
		observer: &defaultObserver{},
		usePool:  true,
		capacity: DefaultConnectionsCapacity,
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

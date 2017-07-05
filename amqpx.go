// Package amqpx provides extentions for Go AMQP library.
//
// It use https://github.com/streadway/amqp as dependency.
package amqpx

import (
	"fmt"

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

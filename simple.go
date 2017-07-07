package amqpx

import (
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// WithoutConnectionsPool will configure a client without a connections pool.
func WithoutConnectionsPool() Option {
	return option(func(options *clientOptions) error {
		options.usePool = false
		return nil
	})
}

// Simple implements the Client interface without a connections pool.
// It will create a new connection every time a channel is requested.
type Simple struct {
	mutex sync.RWMutex

	dialer   Dialer
	observer Observer

	closed bool
}

func newSimple(options *clientOptions) (Client, error) {
	instance := &Simple{
		dialer:   options.dialer,
		observer: options.observer,
	}

	return instance, nil
}

// Channel returns a new amqp's channel from current client unless it's closed.
func (e *Simple) Channel() (*amqp.Channel, error) {
	e.mutex.RLock()
	closed := e.closed
	e.mutex.RUnlock()

	if closed {
		return nil, errors.Wrap(ErrClientClosed, "amqpx: cannot open a new channel")
	}

	connection, err := e.dialer()
	if err != nil {
		return nil, errors.Wrap(err, "amqpx: cannot open a new channel")
	}

	channel, err := connection.Channel()
	if err != nil {
		e.close(connection)
		e.close(channel)
		return nil, errors.Wrap(err, "amqpx: cannot open a new channel")
	}

	return channel, nil
}

// Close closes the client.
func (e *Simple) Close() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.closed = true
	return nil
}

// IsClosed returns if the client is closed.
func (e *Simple) IsClosed() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.closed
}

func (e *Simple) close(connection io.Closer) {
	if connection == nil {
		return
	}

	err := connection.Close()
	if err != nil {
		e.observer.OnClose(err)
	}
}

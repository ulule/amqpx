package amqpx

import (
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Simple implements the Client interface without a connections pool.
// It will use a single connection for multiple channel.
type Simple struct {
	mutex      sync.RWMutex
	dialer     Dialer
	observer   Observer
	connection *amqp.Connection
	closed     bool
	retrier    retrier
}

func newSimple(options *clientOptions) (Client, error) {
	instance := &Simple{
		dialer:   options.dialer,
		observer: options.observer,
		retrier:  newRetrier(options.retry),
	}

	err := instance.newConnection()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// Channel returns a new amqp's channel from current client unless it's closed.
func (e *Simple) Channel() (Channel, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.closed {
		return nil, errors.Wrap(ErrClientClosed, "amqpx: cannot open a new channel")
	}

	ch, err := e.connection.Channel()
	channel := newChannel(ch, e.retrier)
	if err != nil && err != amqp.ErrClosed {
		if ch != nil {
			e.close(channel)
		}
		return nil, errors.Wrap(err, "amqpx: cannot open a new channel")
	}

	// If connection is closed...
	if err == amqp.ErrClosed {
		// Try to open a new one.
		err = e.newConnection()
		if err != nil {
			return nil, err
		}

		// And obtain a new channel.
		ch, err = e.connection.Channel()
		channel := newChannel(ch, e.retrier)
		if err != nil {
			if ch != nil {
				e.close(channel)
			}
			return nil, errors.Wrap(err, "amqpx: cannot open a new channel")
		}
	}

	return channel, nil
}

func (e *Simple) newConnection() error {
	if e.connection != nil {
		e.close(e.connection)
	}

	connection, err := e.dialer.dial(0)
	if err != nil {
		return errors.Wrap(err, "amqpx: cannot open a new connection")
	}

	e.connection = connection
	return nil
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
	err := connection.Close()
	if err != nil {
		e.observer.OnClose(err)
	}
}

var _ Client = (*Simple)(nil)

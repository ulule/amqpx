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
	retriers   *retriers
}

func newSimple(options *clientOptions) (Client, error) {
	instance := &Simple{
		dialer:   options.dialer,
		observer: options.observer,
		retriers: newRetriers(options.retriers),
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
		return nil, errors.Wrap(ErrClientClosed, ErrOpenChannel.Error())
	}

	// Try to acquire a channel.
	channel, err := openChannel(e.connection, e.retriers.channel, e.observer)
	if err != nil {
		return nil, errors.Wrap(err, ErrOpenChannel.Error())
	}

	// If channel is closed, renew the connection and try again.
	if err == amqp.ErrClosed {
		err = e.newConnection()
		if err != nil {
			return nil, err
		}

		channel, err = openChannel(e.connection, e.retriers.channel, e.observer)
		if err != nil {
			return nil, errors.Wrap(err, ErrOpenChannel.Error())
		}
	}

	return channel, nil
}

func (e *Simple) newConnection() error {
	var (
		err        error
		connection *amqp.Connection
	)

	if e.connection != nil {
		e.close(e.connection)
	}

	err = e.retriers.connection.retry(func() error {
		connection, err = e.dialer.dial(0)
		if err != nil {
			return errors.Wrap(err, ErrOpenConnection.Error())
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, ErrRetryExceeded.Error())
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

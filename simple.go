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
		return nil, errors.Wrap(ErrClientClosed, ErrOpenChannel.Error())
	}

	// Try to get a channel.
	channel, err := e.getChannel(e.connection)
	if err != nil {
		return nil, errors.WithStack(ErrOpenChannel)
	}

	// If connection is closed...
	if err == amqp.ErrClosed {
		// Try to open a new one.
		err = e.newConnection()
		if err != nil {
			return nil, err
		}

		channel, err = e.getChannel(e.connection)
		if err != nil {
			return nil, errors.WithStack(ErrOpenChannel)
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
		return errors.Wrap(err, ErrOpenConnection.Error())
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

func (e *Simple) getChannel(conn *amqp.Connection) (Channel, error) {
	var (
		err error
		ch  *amqp.Channel
	)

	err = e.retrier.retry(func() error {
		ch, err = conn.Channel()
		if err != nil && err != amqp.ErrClosed {
			if ch != nil {
				e.close(ch)
			}
			return errors.Wrap(err, ErrOpenChannel.Error())
		}
		return nil
	})

	return newChannel(ch, e.retrier), nil
}

var _ Client = (*Simple)(nil)

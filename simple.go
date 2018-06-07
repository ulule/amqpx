package amqpx

import (
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Simple implements the Client interface without a connections pool.
// It will use a single connection for multiple channel.
type Simple struct {
	mutex        sync.RWMutex
	dialer       Dialer
	observer     Observer
	logger       Logger
	connection   *amqp.Connection
	closed       bool
	retryOptions retriersOptions
}

func newSimple(options *clientOptions) (Client, error) {
	instance := &Simple{
		dialer:       options.dialer,
		observer:     options.observer,
		logger:       options.logger,
		retryOptions: options.retriers,
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
		return nil, errors.Wrap(ErrClientClosed, ErrMessageCannotOpenChannel)
	}

	// Try to acquire a channel.
	channel, err := openChannel(e.connection, e.retryOptions.channel, e.observer, e.logger)
	if err != nil {
		return nil, errors.Wrap(err, ErrMessageCannotOpenChannel)
	}

	// If channel is closed, renew the connection and try again.
	if err == amqp.ErrClosed {
		e.logger.Debug("Channel is closed - opening a new connection...")

		err = e.newConnection()
		if err != nil {
			return nil, err
		}

		channel, err = openChannel(e.connection, e.retryOptions.channel, e.observer, e.logger)
		if err != nil {
			return nil, errors.Wrap(err, ErrMessageCannotOpenChannel)
		}
	}

	return channel, nil
}

func (e *Simple) newConnection() error {
	var (
		err        error
		connection *amqp.Connection
		retry      = newRetrier(e.retryOptions.connection)
	)

	if e.connection != nil {
		e.close(e.connection)
	}

	err = retry.retry(func() error {
		connection, err = e.dialer.dial(0)
		if err != nil {
			e.logger.Error("Failed to obtain a connection - trying again")
			return errors.Wrap(err, ErrMessageCannotOpenConnection)
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, ErrMessageRetryExceeded)
	}

	e.logger.Debug(fmt.Sprintf("Opened connection %s", connection.LocalAddr()))
	e.connection = connection

	return nil
}

// Close closes the client.
func (e *Simple) Close() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.logger.Debug(fmt.Sprintf("Closing connection %s", e.connection.LocalAddr()))
	e.close(e.connection)
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
		e.logger.Error("Failed to close connection")
		e.observer.OnClose(err)
	}

	e.closed = true
}

var _ Client = (*Simple)(nil)

package amqpx

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// A simple client.
type simple struct {
	mutex  sync.RWMutex
	dialer Dialer
	closed bool
}

func NewSimple(dialer Dialer) (Client, error) {
	instance := &simple{
		dialer: dialer,
	}

	return instance, nil
}

// Channel returns a new amqp's channel from current client unless it's closed.
func (e *simple) Channel() (*amqp.Channel, error) {
	e.mutex.RLock()
	closed := e.closed
	e.mutex.RUnlock()

	if closed {
		return nil, errors.Wrap(ErrClientClosed, "amqpx: cannot open a new channel")
	}

	connection, err := e.dialer()
	if err != nil {
		return nil, errors.Wrap(err, "amqpx: cannot open a new connection")
	}

	channel, err := connection.Channel()
	if err != nil {
		if channel != nil {
			// TODO Log/Capture me, but must be silent.
			thr := channel.Close()
			_ = thr
		}
		return nil, errors.Wrap(err, "amqpx: cannot open a new channel")
	}

	return channel, nil
}

// Close closes the client.
func (e *simple) Close() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.closed = true
	return nil
}

// IsClosed returns if the client is closed.
func (e *simple) IsClosed() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.closed
}

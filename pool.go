package amqpx

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Pool implements the Client interface using a connections pool.
// It will reuse a healthy connection from pool when a channel is requested.
type Pool struct {
	mutex       sync.RWMutex
	dialer      Dialer
	observer    Observer
	logger      Logger
	connections []*amqp.Connection
	closed      bool
}

// newPool returns a new client which use a connections pool for amqp's channel.
func newPool(options *clientOptions) (Client, error) {
	instance := &Pool{
		dialer:   options.dialer,
		observer: options.observer,
		logger:   options.logger,
	}

	instance.connections = []*amqp.Connection{}

	for i := 0; i < options.capacity; i++ {
		err := instance.newConnection()
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

// newConnection add a new connection on the connections pool.
func (e *Pool) newConnection() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	idx := len(e.connections)
	connection, err := e.dialer.dial(idx)
	if err != nil {
		e.logger.Error("Failed to obtain a connection")
		return errors.Wrap(err, ErrMessageCannotOpenConnection)
	}

	e.logger.Debug(fmt.Sprintf("Opened connection %s", connection.LocalAddr()))
	e.connections = append(e.connections, connection)
	e.listenOnCloseConnection(idx, connection)

	return nil
}

// releaseConnection remove a connection from the connections pool.
func (e *Pool) releaseConnection(idx int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.connections[idx] = nil
	e.logger.Debug(fmt.Sprintf("Released connection %d", idx))
}

// listenOnCloseConnection will listen on a connection close event.
// If a connection is closed, it will release it from the connections pool and will try to create a new one.
func (e *Pool) listenOnCloseConnection(idx int, connection *amqp.Connection) {
	receiver := make(chan *amqp.Error)
	connection.NotifyClose(receiver)

	go func() {
		err := <-receiver
		if err != nil {
			e.observer.OnClose(err)
		}

		e.releaseConnection(idx)
		e.retryConnection(idx)
	}()
}

// retryConnection will try to open a new connection, unless the client is closed.
// If it succeed, it will add this connection on the connections pool.
func (e *Pool) retryConnection(idx int) {
	e.logger.Debug(fmt.Sprintf("Retrying to open a new connection %d", idx))

	for {
		e.mutex.RLock()
		closed := e.closed
		e.mutex.RUnlock()

		// If client is closed, cancel retry.
		if closed {
			return
		}

		// Try to open a new connection.
		connection, err := e.dialer.dial(idx)
		if err == nil {
			e.mutex.Lock()
			e.logger.Debug(fmt.Sprintf("Opened new connection %s (%d)", connection.LocalAddr(), idx))
			e.connections[idx] = connection
			e.listenOnCloseConnection(idx, connection)
			e.mutex.Unlock()
			return
		}

		// Schedule a retry between 200ms and 1s.
		retry := time.Duration((200 + rand.Intn(801))) * time.Millisecond
		time.Sleep(retry)
	}
}

// Channel returns a new channel from our connections pool.
func (e *Pool) Channel() (*amqp.Channel, error) {
	e.mutex.RLock()
	capacity := len(e.connections)
	offset := rand.Intn(capacity)
	e.mutex.RUnlock()

	for i := 0; i < capacity; i++ {
		idx := (i + offset) % capacity

		e.mutex.RLock()
		connection := e.connections[idx]
		closed := e.closed
		e.mutex.RUnlock()

		if closed {
			return nil, errors.Wrap(ErrClientClosed, ErrMessageCannotOpenChannel)
		}

		if connection != nil {
			channel, err := connection.Channel()
			if err == nil {
				return channel, nil
			}
		}
	}

	return nil, errors.Wrap(ErrNoConnectionAvailable, ErrMessageCannotOpenChannel)
}

// IsClosed returns if the pool is closed.
func (e *Pool) IsClosed() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.closed
}

// Close will closes all remaining connections and marks it as closed.
func (e *Pool) Close() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// If pool is already closed, it's a no-op.
	if e.closed {
		return nil
	}

	e.closed = true
	for i := range e.connections {
		if e.connections[i] != nil {
			conn := e.connections[i]
			e.logger.Debug(fmt.Sprintf("Closing connection %d (%s)", i, conn.LocalAddr()))
			e.close(conn)
		}
	}

	return nil
}

// Length returns connections pool capacity.
func (e *Pool) Length() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return len(e.connections)
}

func (e *Pool) close(connection io.Closer) {
	err := connection.Close()
	if err != nil {
		e.observer.OnClose(err)
	}
}

var _ Client = (*Pool)(nil)

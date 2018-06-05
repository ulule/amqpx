package amqpx

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Pooler implements the Client interface using a connections pool.
// It will reuse a healthy connection from pool when a channel is requested.
type Pooler struct {
	mutex       sync.RWMutex
	dialer      Dialer
	observer    Observer
	connections []*amqp.Connection
	closed      bool
	retrier     retrier
}

// newConnectionsPool returns a new client which use a connections pool for amqp's channel.
func newConnectionsPool(options *clientOptions) (Client, error) {
	// Default channel pool.
	instance := &Pooler{
		dialer:   options.dialer,
		observer: options.observer,
		retrier:  newRetrier(options.retry),
	}

	// Create connections pool.
	instance.connections = []*amqp.Connection{}

	// Open and keep amqp connections.
	for i := 0; i < options.capacity; i++ {
		err := instance.newConnection()
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

// newConnection add a new connection on the connections pool.
func (e *Pooler) newConnection() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	idx := len(e.connections)
	connection, err := e.dialer.dial(idx)
	if err != nil {
		return err
	}

	e.connections = append(e.connections, connection)
	e.listenOnCloseConnection(idx, connection)

	return nil
}

// releaseConnection remove a connection from the connections pool.
func (e *Pooler) releaseConnection(idx int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.connections[idx] = nil
}

// listenOnCloseConnection will listen on a connection close event.
// If a connection is closed, it will release it from the connections pool and will try to create a new one.
func (e *Pooler) listenOnCloseConnection(idx int, connection *amqp.Connection) {
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
func (e *Pooler) retryConnection(idx int) {
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

// Channel returns a new Channel from our connections pool.
func (e *Pooler) Channel() (Channel, error) {
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
			return nil, errors.Wrap(ErrClientClosed, "amqpx: cannot open a new channel")
		}

		if connection != nil {
			return openChannel(connection, e.retrier, e.observer)
		}
	}

	return nil, errors.Wrap(ErrNoConnectionAvailable, "amqpx: cannot open a new channel")
}

// IsClosed returns if the pool is closed.
func (e *Pooler) IsClosed() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.closed
}

// Close will closes all remaining connections and marks it as closed.
func (e *Pooler) Close() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// If pool is already closed, it's a no-op.
	if e.closed {
		return nil
	}

	e.closed = true
	for i := range e.connections {
		if e.connections[i] != nil {
			e.close(e.connections[i])
		}
	}

	return nil
}

// Length returns connections pool capacity.
func (e *Pooler) Length() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return len(e.connections)
}

func (e *Pooler) close(connection io.Closer) {
	err := connection.Close()
	if err != nil {
		e.observer.OnClose(err)
	}
}

var _ Client = (*Pooler)(nil)

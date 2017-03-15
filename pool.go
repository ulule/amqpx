package amqpx

import (
	"errors"
	"io"
	"sync"

	"github.com/streadway/amqp"
)

// Dialer is a function returning an amqp connection
type Dialer func() (*amqp.Connection, error)

// Pooler interface describe a pool implementation
type Pooler interface {

	// Get returns an amqp connection from the pool
	// Closing this connection puts it back to the pool
	Get() (Connector, error)

	// Close closes the pool and all its connections
	// A closed pool cannot be used again
	Close() error

	// Length counts open connections
	Length() int
}

// channelPool implements the Pooler interface using a channel and a mutex
type channelPool struct {
	mutex sync.RWMutex

	dialer Dialer

	connections chan *amqp.Connection
	capacity    int
	closed      bool
}

// ChannelPoolOption is a function modifying a channelPool configuration
type ChannelPoolOption func(*channelPool) error

const (
	// DefaultConnectionsCapacity is default connections channel capacity
	DefaultConnectionsCapacity = 20
)

var (
	// ErrInvalidChannelCapacity occurs when channel has an invalid capacity
	ErrInvalidChannelCapacity = errors.New("invalid channel pool capacity")

	// ErrChannelPoolClosed occurs when operating on a closed channel pool
	ErrChannelPoolClosed = errors.New("channel pool already closed")
)

// Capacity define the capacity of connections channel
func Capacity(capacity int) ChannelPoolOption {
	return func(c *channelPool) error {

		if capacity <= 0 {
			return ErrInvalidChannelCapacity
		}

		c.capacity = capacity
		return nil
	}
}

// NewChannelPool returns a new channelPool
// it calls our dialer and fill the connections channel until bounds are met
func NewChannelPool(dialer Dialer, options ...ChannelPoolOption) (Pooler, error) {

	// Default channel pool.
	pool := &channelPool{
		dialer:   dialer,
		capacity: DefaultConnectionsCapacity,
	}

	// Applies options.
	for _, option := range options {
		err := option(pool)
		if err != nil {
			return nil, err
		}
	}

	// Create connections chan from given options.
	pool.connections = make(chan *amqp.Connection, pool.capacity)

	// Open and store amqp connections.
	for i := 0; i < pool.capacity; i++ {

		connection, err := dialer()
		if err != nil {
			return nil, err
		}

		pool.connections <- connection

	}

	return pool, nil
}

// release puts back a connection to our pool
func (c *channelPool) release(connection *amqp.Connection) error {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// In case channel was closed but given connection is still active.
	if c.closed {
		c.close(connection)
		return nil
	}

	c.connections <- connection
	return nil
}

// Get returns a connection from the pool.
// However, if no connection are avaible, it will block.
func (c *channelPool) Get() (Connector, error) {

	connection, ok := <-c.connections
	if !ok {
		return nil, ErrChannelPoolClosed
	}

	return NewPoolConnection(c, connection), nil
}

// IsClosed returns if the pool is closed.
func (c *channelPool) IsClosed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.closed
}

// Close will closes all remaining connections and marks it as closed.
func (c *channelPool) Close() error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If pool is already closed, it's a no-op.
	if c.closed {
		return nil
	}

	c.closed = true
	close(c.connections)
	for connection := range c.connections {
		c.close(connection)
	}

	return nil
}

func (c *channelPool) close(connection io.Closer) {
	err := connection.Close()
	if err != nil {
		// TODO Log/Capture me, but must to be silent.
	}
}

// Length returns connections channel length
func (c *channelPool) Length() int {
	return len(c.connections)
}

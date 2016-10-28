package amqpx

import (
	"errors"
	"sync"

	"github.com/streadway/amqp"
)

// Dialer is a function returning an amqp connection
type Dialer func() (*amqp.Connection, error)

// Pooler interface describe a pool implementation
type Pooler interface {
	// Get returns an amqp connection from the pool
	// Closing this connection puts it back to the pool
	Get() (*amqp.Connection, error)

	// Close closes the pool and all its connections
	// A closed pool cannot be used again
	Close() error
}

// channelPool implements the Pooler interface using a channel and a mutex
type channelPool struct {
	mu sync.Mutex

	dialer Dialer

	connections chan *amqp.Connection
	length      int
	capacity    int
}

// channelPoolOption is a function modifying a channelPool configuration
type channelPoolOption func(*channelPool) error

const (
	DefaultConnectionsLength   = 10
	DefaultConnectionsCapacity = 20
)

var (
	// ErrInvalidChannelPoolBounds occurs when channel got bad bounds
	ErrInvalidChannelPoolBounds = errors.New("invalid channel pool bounds")

	// ErrChannelPoolAlreadyClosed occurs when operating on an already closed channelPool
	ErrChannelPoolAlreadyClosed = errors.New("channel pool already closed")
)

// Bounds set our channelPool connections channel length and capacity
func Bounds(length, capacity int) channelPoolOption {
	return func(c *channelPool) error {
		if length < 0 || capacity < length {
			return ErrInvalidChannelPoolBounds
		}

		c.length = length
		c.capacity = capacity

		return nil
	}
}

// NewChannelPool returns a new channelPool
// it calls our dialer and fill the connections channel until bounds are met
func NewChannelPool(dialer Dialer, options ...channelPoolOption) (Pooler, error) {
	// default channelPool
	pool := &channelPool{
		dialer:   dialer,
		length:   DefaultConnectionsLength,
		capacity: DefaultConnectionsCapacity,
	}

	// applies options
	for _, option := range options {
		err := option(pool)
		if err != nil {
			return nil, err
		}
	}

	// create connections chan from given options
	pool.connections = make(chan *amqp.Connection, pool.capacity)

	// open and store amqp connections
	for i := 0; i < pool.length; i++ {
		connection, err := dialer()
		if err != nil {
			return nil, err
		}

		pool.connections <- connection
	}

	return pool, nil
}

// putBack puts back a connection to our channelPool connections channel
func (c *channelPool) putBack(connection *amqp.Connection) error {
	if c.connections == nil {
		return ErrChannelPoolAlreadyClosed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.connections <- connection:
		return nil
	default:
		return connection.Close()
	}
}

// wrap wraps a standard amqp.Connection to a PoolConnection
// allow us to override the amqp.Connection.Close method
// instead of closing it we put it back to channelPool using the putBack method
func (c *channelPool) wrap(connection *amqp.Connection) *amqp.Connection {
	p := &PoolConnection{pool: c}
	p.Connection = connection
	return p
}

// Get returns a connection from the channelPool connections channel
// if no connection are avaible, creates a new one calling our dialer
func (c *channelPool) Get() (*amqp.Connection, error) {
	if c.connections == nil {
		return nil, ErrChannelPoolAlreadyClosed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case connection := <-c.connections:
		return c.wrap(connection), nil
	default:
		connection, err := c.dialer()
		if err != nil {
			return nil, err
		}

		return c.wrap(connection), nil
	}
}

// Close closes all channelPool connections and marks it as closed
func (c *channelPool) Close() error {
	if c.connections == nil {
		return ErrChannelPoolAlreadyClosed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for connection := range c.connections {
		connection.Close()
	}
	close(c.connections)
	c.connections = nil

	return nil
}

package amqpx

import "github.com/streadway/amqp"

// Connector interface describe an amqp Connection
type Connector interface {
	// Channel return a new Channel
	Channel() (*amqp.Channel, error)

	// Close closes an amqp Connection
	Close() error
}

// PoolConnnection wrap an amqp.Connection
// It implements the Connector interface
// At close, it puts back its connection to its pool
type PoolConnection struct {
	pool       *channelPool
	connection *amqp.Connection
}

// NewPoolConnection returns a new PoolConnection
func NewPoolConnection(pool *channelPool, connection *amqp.Connection) *PoolConnection {
	return &PoolConnection{
		pool:       pool,
		connection: connection,
	}
}

// Channel returns a new Channel
func (p *PoolConnection) Channel() (*amqp.Channel, error) {
	return p.connection.Channel()
}

// Close puts back the connection to the channelPool connections channel
func (p *PoolConnection) Close() error {
	return p.pool.putBack(p.connection)
}

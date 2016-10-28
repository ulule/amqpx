package amqpx

import "github.com/streadway/amqp"

// PoolConnnection wrap an amqp.Connection
// used to override default Close method
type PoolConnection struct {
	*amqp.Connection

	pool *channelPool
}

// Close puts back the connection to the channelPool connections channel
func (p *PoolConnection) Close() error {
	return p.pool.putBack(p.Connection)
}

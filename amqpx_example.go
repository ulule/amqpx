package amqpx

import "github.com/streadway/amqp"

func NewChannelPoolExample() {
	uri := "amqp://..."

	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(uri)
	}

	// pool will contain 20 amqp connections
	pool, err := NewChannelPool(dialer, Bounds(20, 40))
	if err != nil {
		panic(err)
	}

	connection, err := pool.Get()
	if err != nil {
		panic(err)
	}

	connection.Channel()
	// ...
}
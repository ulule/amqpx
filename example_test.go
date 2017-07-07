package amqpx_test

import (
	"github.com/streadway/amqp"
	"github.com/ulule/amqpx"
)

func ExampleNew() {
	uri := "amqp://..."

	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(uri)
	}

	// This client will contain 20 amqp connections.
	client, err := amqpx.New(dialer, amqpx.WithCapacity(20))
	if err != nil {
		// Handle error...
	}

	channel, err := client.Channel()
	if err != nil {
		// Handle error...
	}

	// ...
	_ = channel
}

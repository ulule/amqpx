package amqpx_test

import (
	"github.com/ulule/amqpx"
)

func ExampleNew() {
	uri := "amqp://..."

	dialer, err := amqpx.SimpleDialer(uri)
	if err != nil {
		// Handle error...
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

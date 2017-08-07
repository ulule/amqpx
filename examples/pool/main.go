package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ulule/amqpx"
)

func main() {
	uri := os.Getenv("AMQP_URI")

	// This dialer will create new connections on a single broker.
	dialer, err := amqpx.SimpleDialer(uri)
	if err != nil {
		panic(err)
	}

	// This client will contain 20 amqp connections.
	client, err := amqpx.New(dialer, amqpx.WithCapacity(20))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Millisecond * 500)

	go func(client amqpx.Client) {
		for {
			channel, err := client.Channel()
			if err != nil {
				fmt.Printf("Error handled %v\n", err)
			} else {
				fmt.Printf("Retrieve AMQP channel %v\n", channel)
			}

			// ...
			_ = channel

			time.Sleep(time.Second)
		}
	}(client)

	time.Sleep(time.Second * 60)
	ticker.Stop()
}
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ulule/amqpx"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var (
		uri      string
		pool     bool
		capacity int
	)

	flag.StringVar(&uri, "uri", "amqp://guest:guest@127.0.0.1:5672/amqpx", "Broker URI")
	flag.BoolVar(&pool, "pool", false, "Enable connections pool")
	flag.IntVar(&capacity, "pool.capacity", amqpx.DefaultConnectionsCapacity, "Number of connections")
	flag.Parse()

	opts := options{
		uri:      uri,
		pool:     pool,
		capacity: capacity,
	}

	client, err := createClient(opts)
	if err != nil {
		return err
	}

	runClient(client)

	return nil
}

type options struct {
	uri      string
	pool     bool
	capacity int
}

func createClient(opts options) (amqpx.Client, error) {
	// This dialer will create new connections on a single broker.
	dialer, err := amqpx.SimpleDialer(opts.uri)
	if err != nil {
		return nil, err
	}

	// By defaults, the client uses of pool of 10 connections.
	// But we only need a single one, we can just pass WithoutConnectionsPool().
	if !opts.pool {
		return amqpx.New(dialer, amqpx.WithoutConnectionsPool())
	}

	// Otherwise, if we need more connections, we can define the capacity
	// with WithCapacity() option.
	return amqpx.New(dialer, amqpx.WithCapacity(opts.capacity))
}

func runClient(client amqpx.Client) {
	var (
		ticker    = time.NewTicker(500 * time.Millisecond)
		interrupt = make(chan os.Signal, 1)
	)

	signal.Notify(interrupt, os.Interrupt)

	defer func() {
		fmt.Printf("Stopping ticker\n")
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			channel, err := client.Channel()
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Printf("Create channel\n")
			err = channel.Close()
			if err != nil {
				log.Fatal(err)
				return
			}
		case <-interrupt:
			return
		}
	}
}

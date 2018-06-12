package main

import (
	"flag"
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
		uri         string
		pool        bool
		capacity    int
		loggerLevel string
	)

	flag.StringVar(&uri, "uri", "amqp://guest:guest@127.0.0.1:5672/amqpx", "Broker URI")
	flag.BoolVar(&pool, "pool", false, "Enable connections pool")
	flag.IntVar(&capacity, "pool.capacity", amqpx.DefaultConnectionsCapacity, "Number of connections")
	flag.StringVar(&loggerLevel, "logger.level", amqpx.LoggerLevelDebugStr, "Logger level")
	flag.Parse()

	if capacity != 0 {
		pool = true
	}

	opts := options{
		uri:         uri,
		pool:        pool,
		capacity:    capacity,
		loggerLevel: amqpx.LoggerLevelFromString(loggerLevel),
	}

	client, err := createClient(opts)
	if err != nil {
		return err
	}

	return runClient(client)
}

type options struct {
	uri         string
	pool        bool
	capacity    int
	loggerLevel amqpx.LoggerLevel
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
		return amqpx.New(dialer,
			amqpx.WithLoggerLevel(opts.loggerLevel),
			amqpx.WithoutConnectionsPool(),
		)
	}

	// Otherwise, if we need more connections, we can define the capacity
	// with WithCapacity() option.
	return amqpx.New(dialer,
		amqpx.WithLoggerLevel(opts.loggerLevel),
		amqpx.WithCapacity(opts.capacity),
	)
}

func runClient(client amqpx.Client) error {
	var (
		ticker    = time.NewTicker(500 * time.Millisecond)
		interrupt = make(chan os.Signal, 1)
	)

	signal.Notify(interrupt, os.Interrupt)

	defer func() {
		ticker.Stop()
		if !client.IsClosed() {
			thr := client.Close()
			_ = thr
		}
	}()

	for {
		select {
		case <-ticker.C:
			channel, err := client.Channel()
			if err != nil {
				return err
			}

			err = channel.Close()
			if err != nil {
				return err
			}

		case <-interrupt:
			return nil
		}
	}
}

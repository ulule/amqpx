# AMQPX

[![CircleCI][circle-img]][circle-url]
[![Documentation][godoc-img]][godoc-url]
![License][license-img]

*An extension for Go AMQP library.*

## Introduction

Amqpx is an extension around [github.com/streadway/amqp](https://github.com/streadway/amqp) to provides high level functionality such as:

 * Client with connection recovery.
 * Client with connections pool.
 * Dialer with cluster support.

## Installation

Using [dep](https://github.com/golang/dep)

```console
dep ensure -add github.com/ulule/amqpx@v2.1.0
```

or `go get`

```console
go get -u github.com/ulule/amqpx
```

## Usage

Amqpx helps you retrieve channels from a `Client`:

```go
// Client interface describes a amqp client.
type Client interface {
	// Channel returns a new Channel from current client unless it's closed.
	Channel() (*amqp.Channel, error)

	// Close closes the client.
	Close() error

	// IsClosed returns if the client is closed.
	IsClosed() bool
}
```

Once a channel is acquired, use it as a simple AMQP channel: **there is no abstractions**.

However, if your channel is closed because there is network connectivity issues on your server _(for example)_, just recycle your channel by querying a new one from your `Client` instance. Depending on your configuration, the client will try to create a new channel from a healthy connection.

### Example

#### Simple

```go
package main

import (
	"github.com/streadway/amqp"
	"github.com/ulule/amqpx"
)

func main() {
	uri := "amqp://guest:guest@127.0.0.1:5672/amqpx"

	dialer, err := amqpx.SimpleDialer(uri)
	if err != nil {
		// Handle error...
	}

	client, err := amqpx.New(dialer, amqpx.WithCapacity(20))
	if err != nil {
		// Handle error...
	}

	channel, err := client.Channel()
	if err != nil {
		// Handle error...
	}

	// Publish and/or receive messages using your channel.

}
```

See [examples](examples/simple) directory for further information.

#### Cluster

```go
package main

import (
	"time"

	"github.com/streadway/amqp"
	"github.com/ulule/amqpx"
)

func main() {
	uris := []string{
		"amqp://guest:guest@127.0.0.1:5672/amqpx",
		"amqp://guest:guest@127.0.0.1:5673/amqpx",
		"amqp://guest:guest@127.0.0.1:5674/amqpx",
	}

	dialer, err := amqpx.ClusterDialer(uris,
		amqpx.WithDialerTimeout(1 * time.Second),
		amqpx.WithDialerTimeout(200 * time.Millisecond),
	)
	if err != nil {
		// Handle error...
	}

	client, err := amqpx.New(dialer, amqpx.WithCapacity(60))
	if err != nil {
		// Handle error...
	}

	channel, err := client.Channel()
	if err != nil {
		// Handle error...
	}

	// Publish and/or receive messages using your channel.

}
```

#### With Observer and Logger

```go
package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/ulule/amqpx"
)

type ZeroLogger struct{}

func (ZeroLogger) Debug(args ...interface{}) {
	log.Debug().Msg(fmt.Sprint(args...))
}

func (ZeroLogger) Info(args ...interface{}) {
	log.Info().Msg(fmt.Sprint(args...))
}

func (ZeroLogger) Warn(args ...interface{}) {
	log.Warn().Msg(fmt.Sprint(args...))
}

func (ZeroLogger) Error(args ...interface{}) {
	log.Error().Msg(fmt.Sprint(args...))
}

type ZeroObserver struct{}

func (ZeroObserver) OnError(err error) {
	log.Error().Err(err).Msg("An error has occurred")
}

func (ZeroObserver) OnClose(err error) {
	log.Warn().Err(err).Msg("A close/shutdown error occured")
}

func main() {
	uris := []string{
		"amqp://guest:guest@127.0.0.1:5672/amqpx",
		"amqp://guest:guest@127.0.0.1:5673/amqpx",
		"amqp://guest:guest@127.0.0.1:5674/amqpx",
	}

	dialer, err := amqpx.ClusterDialer(uris,
		amqpx.WithDialerTimeout(1 * time.Second),
		amqpx.WithDialerTimeout(200 * time.Millisecond),
	)
	if err != nil {
		// Handle error...
	}

	client, err := amqpx.New(dialer,
		amqpx.WithCapacity(60),
		amqpx.WithLogger(ZeroLogger{}),
		amqpx.WithObserver(ZeroObserver{}),
	)
	if err != nil {
		// Handle error...
	}

	channel, err := client.Channel()
	if err != nil {
		// Handle error...
	}

	// Publish and/or receive messages using your channel.

}
```

## License

This is Free Software, released under the [`MIT License`][license-url].

## Contributing

* Fix [bugs](https://github.com/ulule/amqpx/issues)
* Fork the [project](https://github.com/ulule/amqpx)
* Submit new [idea](https://github.com/ulule/amqpx/issues)
* Ping us on twitter:
  * [@novln_](https://twitter.com/novln_)
  * [@oibafsellig](https://twitter.com/oibafsellig)
  * [@thoas](https://twitter.com/thoas)

**Don't hesitate ;)**

[godoc-url]: https://godoc.org/github.com/ulule/amqpx
[godoc-img]: https://godoc.org/github.com/ulule/amqpx?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: LICENSE
[circle-url]: https://circleci.com/gh/ulule/amqpx/tree/master
[circle-img]: https://circleci.com/gh/ulule/amqpx.svg?style=shield&circle-token=a76e635936a3dc466d8ee83d9c03524598bae4b8

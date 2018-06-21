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
go get -u github.com/ulule/loukoum
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

```go
package foobar

import (
	"github.com/streadway/amqp"
	"github.com/ulule/amqpx"
)

func main() {
	uri := "amqp://guest:guest@127.0.0.1:5672/"

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

See [examples](examples/simple) directory for more information.


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

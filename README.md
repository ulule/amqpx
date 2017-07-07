# AMQPX

An extension for Go AMQP library.

Currently supports:

 * Client with connection recovery.
 * Client with connections pool.

## Installation

```
$ go get -u github.com/ulule/amqpx
```

## Usage

```go
package foobar

import (
	"github.com/streadway/amqp"
	"github.com/ulule/amqpx"
)

func main() {
	uri := "amqp://guest:guest@127.0.0.1:5672/"

	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(uri)
	}

	client, err := amqpx.New(dialer, amqpx.WithCapacity(20))
	if err != nil {
		// Handle error...
	}

	channel, err := client.Channel()
	if err != nil {
		// Handle error...
	}

	// ...

}
```

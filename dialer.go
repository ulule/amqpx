package amqpx

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var (
	// ErrBrokerURIRequired occurs when a dialer has no broker URI.
	ErrBrokerURIRequired = fmt.Errorf("broker URI is required")
)

// Dialer is an interface that return a new amqp connection.
// In order to instantiate a new Dialer, please use SimpleDialer or ClusterDialer.
type Dialer interface {
	dial(id int) (*amqp.Connection, error)
}

// SimpleDialer gives a Dialer that use a simple broker.
func SimpleDialer(uri string) (Dialer, error) {
	if uri == "" {
		return nil, errors.Wrap(ErrBrokerURIRequired, "amqpx: cannot create a new dialer")
	}

	dialer := &simpleDialer{uri: uri}

	return dialer, nil
}

// ClusterDialer is a Dialer that use a cluster of broker.
func ClusterDialer(list []string) (Dialer, error) {
	if len(list) == 0 {
		return nil, errors.Wrap(ErrBrokerURIRequired, "amqpx: cannot create a new dialer")
	}

	dialer := &clusterDialer{list: list}
	return dialer, nil
}

type simpleDialer struct {
	uri string
}

func (e *simpleDialer) dial(id int) (*amqp.Connection, error) {
	return amqp.Dial(e.uri)
}

type clusterDialer struct {
	list []string
}

func (e *clusterDialer) dial(id int) (*amqp.Connection, error) {
	idx := (id) % len(e.list)
	uri := e.list[idx]
	return amqp.Dial(uri)
}

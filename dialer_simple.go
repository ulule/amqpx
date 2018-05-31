package amqpx

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// SimpleDialer gives a Dialer that use a simple broker.
func SimpleDialer(uri string, options ...DialerOptions) (Dialer, error) {
	if uri == "" {
		return nil, errors.Wrap(
			ErrBrokerURIRequired,
			"amqpx: cannot create a new dialer")
	}

	return &simpleDialer{
		DialerOptions: NewDialerOptions(options...),
		uri:           uri,
	}, nil
}

type simpleDialer struct {
	DialerOptions
	uri string
}

func (e *simpleDialer) dial(id int) (*amqp.Connection, error) {
	return amqp.DialConfig(e.uri, amqp.Config{
		Dial:      dialer(e.Timeout),
		Heartbeat: e.Heartbeat,
	})
}

var _ Dialer = (*simpleDialer)(nil)

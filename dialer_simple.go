package amqpx

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// SimpleDialer gives a Dialer that uses a simple broker.
func SimpleDialer(uri string, options ...DialerOption) (Dialer, error) {
	if uri == "" {
		return nil, errors.Wrap(ErrBrokerURIRequired, ErrMessageCannotCreateDialer)
	}

	opts := newDialerOptions()
	for _, option := range options {
		err := option.apply(&opts)
		if err != nil {
			return nil, errors.Wrap(err, ErrMessageCannotCreateDialer)
		}
	}

	return &simpleDialer{
		dialerOptions: opts,
		uri:           uri,
	}, nil
}

type simpleDialer struct {
	dialerOptions
	uri string
}

// Timeout implements Dialer interface.
func (e simpleDialer) Timeout() time.Duration {
	return e.timeout
}

// Heartbeat implements Dialer interface.
func (e simpleDialer) Heartbeat() time.Duration {
	return e.heartbeat
}

// URLs implements Dialer interface.
func (e simpleDialer) URLs() []string {
	return []string{e.uri}
}

func (e simpleDialer) dial(id int) (*amqp.Connection, error) {
	return amqp.DialConfig(e.uri, amqp.Config{
		Dial:      dialer(e.timeout),
		Heartbeat: e.heartbeat,
	})
}

var _ Dialer = (*simpleDialer)(nil)

package amqpx

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ClusterDialer is a Dialer that use a cluster of broker.
func ClusterDialer(list []string, options ...DialerOption) (Dialer, error) {
	if len(list) == 0 {
		return nil, errors.Wrap(
			ErrBrokerURIRequired,
			"amqpx: cannot create a new dialer")
	}

	opts := newDialerOptions()
	for _, option := range options {
		err := option.apply(&opts)
		if err != nil {
			return nil, errors.Wrap(err, "amqpx: cannot create a new dialer")
		}
	}

	return &clusterDialer{
		dialerOptions: opts,
		list:          list,
	}, nil
}

type clusterDialer struct {
	dialerOptions
	list []string
}

// Timeout implements Dialer interface.
func (e clusterDialer) Timeout() time.Duration {
	return e.timeout
}

// Heartbeat implements Dialer interface.
func (e clusterDialer) Heartbeat() time.Duration {
	return e.heartbeat
}

// URLs implements Dialer interface.
func (e clusterDialer) URLs() []string {
	return e.list
}

// dial implements Dialer interface.
func (e clusterDialer) dial(id int) (*amqp.Connection, error) {
	idx := (id) % len(e.list)
	uri := e.list[idx]
	return amqp.DialConfig(uri, amqp.Config{
		Dial:      dialer(e.timeout),
		Heartbeat: e.heartbeat,
	})
}

var _ Dialer = (*clusterDialer)(nil)

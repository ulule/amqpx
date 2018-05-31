package amqpx

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ClusterDialer is a Dialer that use a cluster of broker.
func ClusterDialer(list []string, options ...DialerOptions) (Dialer, error) {
	if len(list) == 0 {
		return nil, errors.Wrap(
			ErrBrokerURIRequired,
			"amqpx: cannot create a new dialer")
	}

	return &clusterDialer{
		DialerOptions: NewDialerOptions(options...),
		list:          list,
	}, nil
}

type clusterDialer struct {
	DialerOptions
	list []string
}

func (e *clusterDialer) dial(id int) (*amqp.Connection, error) {
	idx := (id) % len(e.list)
	uri := e.list[idx]
	return amqp.DialConfig(uri, amqp.Config{
		Dial:      dialer(e.Timeout),
		Heartbeat: e.Heartbeat,
	})
}

var _ Dialer = (*clusterDialer)(nil)

package amqpx

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Dialer is an interface that returns a new amqp connection.
// In order to instantiate a new Dialer, please use SimpleDialer or ClusterDialer.
type Dialer interface {
	Timeout() time.Duration
	Heartbeat() time.Duration
	URLs() []string
	dial(id int) (*amqp.Connection, error)
}

func dialer(timeout time.Duration) func(network string, address string) (net.Conn, error) {
	return func(network string, address string) (net.Conn, error) {
		// Dial a remote address with a timeout.
		conn, err := net.DialTimeout(network, address, timeout)
		if err != nil {
			// TODO Better error message
			return nil, errors.Wrap(err, "dial has timeout")
		}

		// Heartbeating hasn't started yet, don't stall forever on a dead server.
		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			// TODO Better error message
			return nil, errors.Wrap(err, "cannot define a read timeout")
		}

		return conn, nil
	}
}

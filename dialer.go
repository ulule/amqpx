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
	dial(id int) (*amqp.Connection, error)
}

func dialer(timeout time.Duration) func(network string, address string) (net.Conn, error) {
	return func(network string, address string) (net.Conn, error) {

		// Dial a remote address with a timeout.
		conn, err := net.DialTimeout(network, address, timeout)
		if err != nil {
			return nil, errors.Wrap(err, ErrMessageDialTimeout)
		}

		// Heartbeating hasn't started yet, don't stall forever to receive packets from server.
		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return nil, errors.Wrap(err, ErrMessageReadTimeout)
		}

		// Also, don't stall forever when sending packets to server.
		err = conn.SetWriteDeadline(time.Now().Add(timeout))
		if err != nil {
			return nil, errors.Wrap(err, ErrMessageWriteTimeout)
		}

		return conn, nil
	}
}

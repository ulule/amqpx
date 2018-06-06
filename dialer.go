package amqpx

import (
	"net"
	"time"

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

type amqpDialer func(network string, addr string) (net.Conn, error)

func dialer(timeout time.Duration) amqpDialer {
	return func(network string, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			return nil, err
		}

		// Heartbeating hasn't started yet, don't stall forever on a dead server.
		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			return nil, err
		}

		return conn, nil
	}
}

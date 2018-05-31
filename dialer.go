package amqpx

import (
	"net"
	"time"

	"github.com/streadway/amqp"
)

// Dialer is an interface that return a new amqp connection.
// In order to instantiate a new Dialer, please use SimpleDialer or ClusterDialer.
type Dialer interface {
	dial(id int) (*amqp.Connection, error)
}

// DialerOptions are options given to dialer instance.
type DialerOptions struct {
	Timeout   time.Duration
	Heartbeat time.Duration
}

// NewDialerOptions returns a new DialerOptions instance.
func NewDialerOptions(fromOptions ...DialerOptions) DialerOptions {
	dialerOptions := DialerOptions{
		Timeout:   defaultDialerTimeout,
		Heartbeat: defaultDialerHeartbeat,
	}

	if len(fromOptions) > 0 {
		opts := fromOptions[0]

		if opts.Timeout != 0 {
			dialerOptions.Timeout = opts.Timeout
		}

		if opts.Heartbeat != 0 {
			dialerOptions.Heartbeat = opts.Heartbeat
		}
	}

	return dialerOptions
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

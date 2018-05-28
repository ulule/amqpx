package amqpx

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Defaults
var (
	defaultDialerTimeout   = 30 * time.Second
	defaultDialerHeartbeat = 10 * time.Second
)

// Errors
var (
	// ErrBrokerURIRequired occurs when a dialer has no broker URI.
	ErrBrokerURIRequired = fmt.Errorf("broker URI is required")
)

// Dialer is an interface that return a new amqp connection.
// In order to instantiate a new Dialer, please use SimpleDialer or ClusterDialer.
type Dialer interface {
	dial(id int) (*amqp.Connection, error)
}

// -----------------------------------------------------------------------------
// Options
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Simple
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Cluster
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Custom dialer
// -----------------------------------------------------------------------------

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

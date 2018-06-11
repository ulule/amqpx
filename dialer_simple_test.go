package amqpx_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestDialer_Simple(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer(brokerURI)
	is.NoError(err)
	is.NotNil(dialer)
}

func TestDialer_Simple_URIRequired(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer("")
	is.Error(err)
	is.Nil(dialer)
	is.Equal(amqpx.ErrBrokerURIRequired, errors.Cause(err))
}

func TestDialer_Simple_WithDialerTimeout(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer(brokerURI)
	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(amqpx.DefaultDialerTimeout, dialer.Timeout())

	dialer, err = amqpx.SimpleDialer(brokerURI, amqpx.WithDialerTimeout(dialerTimeout))
	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(dialerTimeout, dialer.Timeout())
}

func TestDialer_Simple_WithDialerHeartbeat(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer(brokerURI)
	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(amqpx.DefaultDialerHeartbeat, dialer.Heartbeat())

	dialer, err = amqpx.SimpleDialer(brokerURI, amqpx.WithDialerHeartbeat(dialerHeartbeat))
	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(dialerHeartbeat, dialer.Heartbeat())
}

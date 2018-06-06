package amqpx_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestDialer_Cluster(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.ClusterDialer(brokerURIs)
	is.NoError(err)
	is.NotNil(dialer)
}

func TestDialer_Cluster_URIRequired(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.ClusterDialer(nil)
	is.Error(err)
	is.Nil(dialer)
	is.Equal(amqpx.ErrBrokerURIRequired, errors.Cause(err))

	dialer, err = amqpx.ClusterDialer([]string{})
	is.Error(err)
	is.Nil(dialer)
	is.Equal(amqpx.ErrBrokerURIRequired, errors.Cause(err))
}

func TestDialer_Cluster_WithDialerTimeout(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.ClusterDialer(brokerURIs)
	is.NoError(err)
	is.NotNil(dialer)

	dialer, err = amqpx.ClusterDialer(
		brokerURIs,
		amqpx.WithDialerTimeout(dialerTimeout))

	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(dialerTimeout, dialer.Timeout())
}

func TestDialer_Cluster_WithDialerHeartbeat(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.ClusterDialer(brokerURIs)
	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(amqpx.DefaultDialerHeartbeat, dialer.Heartbeat())

	dialer, err = amqpx.ClusterDialer(
		brokerURIs,
		amqpx.WithDialerHeartbeat(dialerHeartbeat))

	is.NoError(err)
	is.NotNil(dialer)
	is.Equal(dialerHeartbeat, dialer.Heartbeat())
}

package amqpx_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestDialer_Cluster(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.ClusterDialer(nil)
	is.Error(err)
	is.Nil(dialer)
	is.Equal(amqpx.ErrBrokerURIRequired, errors.Cause(err))

	dialer, err = amqpx.ClusterDialer([]string{})
	is.Error(err)
	is.Nil(dialer)
	is.Equal(amqpx.ErrBrokerURIRequired, errors.Cause(err))

	dialer, err = amqpx.ClusterDialer([]string{brokerURI})
	is.NoError(err)
	is.NotNil(dialer)
}

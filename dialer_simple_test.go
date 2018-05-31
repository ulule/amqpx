package amqpx_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestDialer_Simple(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer("")
	is.Error(err)
	is.Nil(dialer)
	is.Equal(amqpx.ErrBrokerURIRequired, errors.Cause(err))

	dialer, err = amqpx.SimpleDialer(brokerURI)
	is.NoError(err)
	is.NotNil(dialer)
}

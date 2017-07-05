package amqpx

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func newSimpleClient() (Client, error) {
	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(brokerURI)
	}

	return NewSimple(dialer)
}

func TestSimpleClient_Channel(t *testing.T) {
	is := require.New(t)

	client, err := newSimpleClient()
	is.NoError(err)
	is.NotNil(client)
	defer client.Close()

	channel, err := client.Channel()
	is.NoError(err)
	is.NotNil(channel)

	err = channel.Close()
	is.NoError(err)
}

func TestSimpleClient_Close(t *testing.T) {
	is := require.New(t)

	client, err := newSimpleClient()
	is.NoError(err)
	is.NotNil(client)
	defer client.Close()

	is.False(client.IsClosed())
	is.NoError(client.Close())
	is.True(client.IsClosed())
}

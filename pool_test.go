package amqpx_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ulule/amqpx"
)

func TestPoolClient(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithCapacity(7))
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Pool{}, client)

	pool := client.(*amqpx.Pool)
	is.Equal(7, pool.Length())

	client, err = NewClient()
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Pool{}, client)

	pool = client.(*amqpx.Pool)
	is.Equal(amqpx.DefaultConnectionsCapacity, pool.Length())
}

func TestPoolClient_WithInvalidBrokerURI(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer(invalidBrokerURI)
	is.NoError(err)

	capacity := 1
	client, err := amqpx.New(dialer, amqpx.WithCapacity(capacity))
	is.Error(err)
	is.Nil(client)
	is.Contains(err.Error(), amqpx.ErrMessageDialTimeout)
	is.Contains(err.Error(), amqpx.ErrMessageCannotOpenConnection)
}

func TestPoolClient_Channel(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient()
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	channel, err := client.Channel()
	is.NoError(err)
	is.NotNil(channel)
	is.NoError(channel.Close())
}

func TestPoolClient_Close(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient()
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	is.False(client.IsClosed())
	is.NoError(client.Close())
	is.True(client.IsClosed())
}

func TestPoolClient_ConcurrentAccess(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient()
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < numberOfConcurrentAccess; i++ {
		wg.Add(1)
		go testClientConcurrentAccess(is, client, wg)
	}

	// Wait a few seconds to simulate a SIGKILL
	time.Sleep(sigkillSleep)
	is.NoError(client.Close())
	wg.Wait()
}

package amqpx_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ulule/amqpx"
)

func TestSimpleClient(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithoutConnectionsPool())
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Simple{}, client)
}

func TestSimpleClient_WithInvalidBrokerURI(t *testing.T) {
	is := NewRunner(t)

	dialer, err := amqpx.SimpleDialer(invalidBrokerURI)
	is.NoError(err)

	client, err := amqpx.New(dialer, amqpx.WithoutConnectionsPool())
	is.Error(err)
	is.Nil(client)
	is.Contains(err.Error(), amqpx.ErrMessageDialTimeout)
	is.Contains(err.Error(), amqpx.ErrMessageCannotOpenConnection)
}

func TestSimpleClient_Close(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithoutConnectionsPool())
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	is.False(client.IsClosed())
	is.NoError(client.Close())
	is.True(client.IsClosed())
}

func TestSimpleClient_ConcurrentAccess(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithoutConnectionsPool())
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

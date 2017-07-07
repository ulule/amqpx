package amqpx_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestSimpleClient_Client(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithoutConnectionsPool())
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Simple{}, client)
}

func TestSimpleClient_Channel(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithoutConnectionsPool())
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	channel, err := client.Channel()
	is.NoError(err)
	is.NotNil(channel)

	err = channel.Close()
	is.NoError(err)
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
	wg := &sync.WaitGroup{}

	client, err := NewClient(amqpx.WithoutConnectionsPool())
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {

			time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)

			channel, err := client.Channel()
			wg.Done()

			defer func() {
				if channel != nil {
					thr := channel.Close()
					_ = thr
				}
			}()

			if err != nil && errors.Cause(err) != amqpx.ErrClientClosed {
				is.NoError(err)
			}

		}()
	}

	// Wait a few seconds to simulate a SIGKILL
	time.Sleep(2 * time.Second)
	err = client.Close()
	is.NoError(err)

	wg.Wait()
}

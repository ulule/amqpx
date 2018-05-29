package amqpx_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/ulule/amqpx"
)

func TestPoolerClient_Client(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithCapacity(7))
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Pooler{}, client)
	pooler := client.(*amqpx.Pooler)
	is.Equal(7, pooler.Length())

	client, err = NewClient()
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Pooler{}, client)
	pooler = client.(*amqpx.Pooler)
	is.Equal(amqpx.DefaultConnectionsCapacity, pooler.Length())
}

func TestPoolerClient_Channel(t *testing.T) {
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

	err = channel.Close()
	is.NoError(err)
}

func TestPoolerClient_Close(t *testing.T) {
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

func TestPoolerClient_ConcurrentAccess(t *testing.T) {
	is := NewRunner(t)
	wg := &sync.WaitGroup{}

	client, err := NewClient()
	is.NoError(err)
	is.NotNil(client)
	defer func() {
		is.NoError(client.Close())
	}()

	for i := 0; i < 8000; i++ {
		wg.Add(1)
		go func() {

			time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)

			var channel *amqp.Channel
			channel, err = client.Channel()
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

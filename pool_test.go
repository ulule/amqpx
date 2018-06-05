package amqpx_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestPoolClient_Client(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithCapacity(7))
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Pool{}, client)

	pooler := client.(*amqpx.Pool)
	is.Equal(7, pooler.Length())

	client, err = NewClient()
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Pool{}, client)

	pooler = client.(*amqpx.Pool)
	is.Equal(amqpx.DefaultConnectionsCapacity, pooler.Length())
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

	err = channel.Close()
	is.NoError(err)
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
	is.NotNil(client)
	is.NoError(err)
	defer func() {
		is.NoError(client.Close())
	}()

	var wg sync.WaitGroup
	for i := 0; i < concurrentAccessNChannels; i++ {
		wg.Add(1)

		go func(clt amqpx.Client, w *sync.WaitGroup) {
			time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)

			ch, cherr := clt.Channel()
			w.Done()
			defer func() {
				if ch != nil {
					_ = ch.Close()
				}
			}()

			if cherr != nil && errors.Cause(cherr) != amqpx.ErrClientClosed {
				is.NoError(cherr)
			}
		}(client, &wg)
	}

	// Wait a few seconds to simulate a SIGKILL
	time.Sleep(sigkillSleep)
	is.NoError(client.Close())
	wg.Wait()
}

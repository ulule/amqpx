package amqpx_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/ulule/amqpx"
)

func TestSimpleClient(t *testing.T) {
	is := NewRunner(t)

	client, err := NewClient(amqpx.WithoutConnectionsPool())
	is.NoError(err)
	is.NotNil(client)
	is.IsType(&amqpx.Simple{}, client)
}

func TestSimpleClient_WithExponentialConnectionRetry(t *testing.T) {
	var (
		is              = NewRunner(t)
		initialInterval = 10 * time.Millisecond
		maxInterval     = 20 * time.Millisecond
		maxElapsedTime  = 100 * time.Millisecond
	)

	dialer, err := amqpx.SimpleDialer(invalidBrokerURI)
	is.NoError(err)

	client, err := amqpx.New(
		dialer,
		amqpx.WithoutConnectionsPool(),
		amqpx.WithExponentialConnectionRetry(
			initialInterval,
			maxInterval,
			maxElapsedTime))

	is.Nil(client)
	is.NotNil(err)
	is.Contains(err.Error(), amqpx.ErrRetryExceeded.Error())
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

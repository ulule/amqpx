package amqpx

import (
	"sync"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func TestChannelPool_Get(t *testing.T) {
	pool, err := newChannelPool()
	if err != nil {
		t.Error(err)
	}
	if pool != nil {
		defer func() {
			err := pool.Close()
			if err != nil {
				t.Error(err)
			}
		}()
	}

	// check length at init
	if pool.Length() != DefaultConnectionsCapacity {
		t.Errorf("Invalid init pool length, expected '%d', got '%d'", DefaultConnectionsCapacity, pool.Length())
	}

	// check length after a get
	connection, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	if pool.Length() != DefaultConnectionsCapacity-1 {
		t.Errorf("Invalid pool length after single get, expected '%d', got '%d'", DefaultConnectionsCapacity-1, pool.Length())
	}

	// check length after a close
	connection.Close()
	if pool.Length() != DefaultConnectionsCapacity {
		t.Errorf("Invalid pool length after close, expected '%d', got '%d'", DefaultConnectionsCapacity, pool.Length())
	}

	// use all connections pool
	connections := make([]Connector, DefaultConnectionsCapacity)
	for i := 0; i < DefaultConnectionsCapacity; i++ {
		connections[i], err = pool.Get()
		if err != nil {
			t.Error(err)
		}
	}

	if pool.Length() != 0 {
		t.Errorf("Invalid pool length after fill, expected '%d', got '%d'", 0, pool.Length())
	}

}

func TestChannelPool_Close(t *testing.T) {
	p, err := newChannelPool()
	if err != nil {
		t.Error(err)
	}

	p.Close()

	pool := p.(*channelPool)

	if pool.closed != true {
		t.Error("Connections should be closed")
	}

	_, err = pool.Get()
	if err != ErrChannelPoolClosed {
		t.Errorf("Get after Close expect error '%s', got '%s'", ErrChannelPoolClosed, err)
	}

	if pool.Length() != 0 {
		t.Errorf("Connection pool should be empty after close")
	}

}

func TestChannelPool_ConcurrentAccess(t *testing.T) {
	wg := &sync.WaitGroup{}

	pool, err := newChannelPool(Capacity(5))
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 8000; i++ {
		wg.Add(1)
		go func(id int) {

			connector, err := pool.Get()
			wg.Done()

			if connector != nil {
				defer func() {
					thr := connector.Close()
					if thr != nil {
						t.Errorf("unexpected error: %s", thr)
					}
				}()
			}

			if err == nil || err == ErrChannelPoolClosed {
				return
			}

			t.Errorf("unexpected error: %s", err)

		}(i)
	}

	// Wait a few seconds to simulate a SIGKILL
	time.Sleep(5 * time.Second)

	err = pool.Close()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

}

func TestChannelPool_Capacity(t *testing.T) {
	scenario := []struct {
		capacity int
		err      error
	}{
		{1, nil},
		{2, nil},
		{30, nil},
		{255, nil},
		{0, ErrInvalidChannelCapacity},
		{-6, ErrInvalidChannelCapacity},
		{-700, ErrInvalidChannelCapacity},
	}

	for _, tt := range scenario {

		p, err := newChannelPool(Capacity(tt.capacity))

		if err != tt.err {
			t.Errorf("unexpected error for %d: %s", tt.capacity, err)
		}

		if tt.err == nil && p == nil {
			t.Errorf("channel pool was expected for %d", tt.capacity)
		}
	}

}

var (
	brokerURI = "amqp://guest:guest@127.0.0.1:5672/"
)

func newChannelPool(options ...ChannelPoolOption) (Pooler, error) {
	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(brokerURI)
	}

	return NewChannelPool(dialer, options...)
}

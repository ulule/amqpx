package amqpx

import (
	"testing"

	"github.com/streadway/amqp"
)

var (
	brokerUri = "amqp://guest:guest@127.0.0.1:5672/"
)

func newChannelPool() (Pooler, error) {
	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(brokerUri)
	}

	return NewChannelPool(dialer)
}

func TestChannelPool_Get(t *testing.T) {
	pool, err := newChannelPool()
	if err != nil {
		t.Error(err)
	}
	defer pool.Close()

	// check length at init
	if pool.Length() != DefaultConnectionsLength {
		t.Errorf("Invalid init pool length, expected '%d', got '%d'", DefaultConnectionsLength, pool.Length())
	}

	// check length after a get
	connection, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	if pool.Length() != DefaultConnectionsLength-1 {
		t.Errorf("Invalid pool length after single get, expected '%d', got '%d'", DefaultConnectionsLength-1, pool.Length())
	}

	// check length after a close
	connection.Close()
	if pool.Length() != DefaultConnectionsLength {
		t.Errorf("Invalid pool length after close, expected '%d', got '%d'", DefaultConnectionsLength, pool.Length())
	}

	// use all connections pool
	connections := make([]Connector, DefaultConnectionsCapacity)
	for i := 0; i < DefaultConnectionsCapacity-DefaultConnectionsLength; i++ {
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

	if pool.connections != nil {
		t.Error("After close, connections pool should be nil")
	}

	_, err = pool.Get()
	if err != ErrChannelPoolAlreadyClosed {
		t.Errorf("Get after Close expect error '%s', got '%s'", ErrChannelPoolAlreadyClosed, err)
	}
}

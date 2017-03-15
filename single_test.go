package amqpx

import (
	"testing"

	"github.com/streadway/amqp"
)

func newSinglePool() (Pooler, error) {
	dialer := func() (*amqp.Connection, error) {
		return amqp.Dial(brokerURI)
	}

	return NewSinglePool(dialer)
}

func TestSinglePool_Get(t *testing.T) {
	pool, err := newChannelPool()
	if err != nil {
		t.Error(err)
	}
	defer pool.Close()

	connection, err := pool.Get()
	if err != nil {
		t.Error(err)
	}

	err = connection.Close()
	if err != nil {
		t.Errorf("Closing a connection should not return an error, got '%s'", err)
	}
}

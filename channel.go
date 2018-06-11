package amqpx

import (
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Channel describes an amqp channel.
type Channel interface {
	//
	// Close initiate a clean channel closure by sending a close message with the error
	// code set to '200'.
	//
	// It is safe to call this method multiple times.
	//
	Close() error
	Qos(prefetchCount int, prefetchSize int, global bool) error
	Cancel(consumer string, noWait bool) error
	QueueDeclare(name string, durable bool, autoDelete bool,
		exclusive bool, noWait bool, args Table) (Queue, error)
	QueueDeclarePassive(name string, durable bool, autoDelete bool,
		exclusive bool, noWait bool, args Table) (Queue, error)
	QueueInspect(name string) (Queue, error)
	QueueBind(name string, key string, exchange string, noWait bool, args Table) error
	QueueUnbind(name string, key string, exchange string, args Table) error
	QueuePurge(name string, noWait bool) (int, error)
	QueueDelete(name string, ifUnused bool, ifEmpty bool, noWait bool) (int, error)
	Consume(queue string, consumer string, autoAck bool, exclusive bool,
		noLocal bool, noWait bool, args Table) (<-chan Delivery, error)
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool,
		internal bool, noWait bool, args Table) error
	ExchangeDeclarePassive(name string, kind string, durable bool, autoDelete bool,
		internal bool, noWait bool, args Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeBind(destination string, key string, source string, noWait bool, args Table) error
	ExchangeUnbind(destination string, key string, source string, noWait bool, args Table) error
	Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error
	Get(queue string, autoAck bool) (msg Delivery, ok bool, err error)
	//
	// Ack acknowledges a delivery by its delivery tag when having been consumed with
	// Channel.Consume or Channel.Get.
	//
	// Ack acknowledges all message received prior to the delivery tag when multiple
	// is true.
	//
	Ack(tag uint64, multiple bool) error
	//
	// Nack negatively acknowledges a delivery by its delivery tag.  Prefer this
	// method to notify the server that you were not able to process this delivery and
	// it must be redelivered or dropped.
	//
	Nack(tag uint64, multiple bool, requeue bool) error
	//
	// Reject negatively acknowledges a delivery by its delivery tag.  Prefer Nack
	// over Reject when communicating with a RabbitMQ server because you can Nack
	// multiple messages, reducing the amount of protocol messages to exchange.
	//
	Reject(tag uint64, requeue bool) error
}

// ChannelWrapper wraps a amqp channel to provide a retry mechanism.
// It implements Channel interface.
type ChannelWrapper struct {
	mutex    sync.RWMutex
	client   Client
	observer Observer
	logger   Logger
	retrier  retrier
	channel  *amqp.Channel
	qos      *channelQos
}

// ChannelWrapper configuration for Qos.
type channelQos struct {
	prefetchCount int
	prefetchSize  int
	global        bool
}

// NewChannelWrapper creates a new Channel.
func NewChannelWrapper(client Client) Channel {
	return &ChannelWrapper{
		client:   client,
		logger:   client.getLogger(),
		observer: client.getObserver(),
		retrier:  client.newRetrier(),
	}
}

// TODO (novln): Decide what and how we obtain a channel on this wrapper.
// TODO (novln): Implement release channel on error.
// TODO (novln): Implement Qos on new channel.

func (ch *ChannelWrapper) Close() error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	if ch.channel != nil {
		err := ch.channel.Close()
		if err != nil {
			ch.logger.Error("Failed to close channel")
			ch.observer.OnClose(err)
			return errors.Wrap(err, ErrMessageCannotCloseChannel)
		}
		ch.channel = nil
	}

	return nil
}

func (ch *ChannelWrapper) Qos(prefetchCount int, prefetchSize int, global bool) error {
	ch.mutex.Lock()
	ch.qos = &channelQos{
		prefetchCount: prefetchCount,
		prefetchSize:  prefetchSize,
		global:        global,
	}
	ch.mutex.Unlock()
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.Qos(prefetchCount, prefetchSize, global)
	})
}

func (ch *ChannelWrapper) Cancel(consumer string, noWait bool) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.Cancel(consumer, noWait)
	})
}

func (ch *ChannelWrapper) QueueDeclare(name string, durable bool, autoDelete bool,
	exclusive bool, noWait bool, args Table) (queue Queue, err error) {
	err = ch.handle(func(channel *amqp.Channel) error {
		queue, err = channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
		return err
	})
	return queue, err
}

func (ch *ChannelWrapper) QueueDeclarePassive(name string, durable bool, autoDelete bool,
	exclusive bool, noWait bool, args Table) (queue Queue, err error) {
	err = ch.handle(func(channel *amqp.Channel) error {
		queue, err = channel.QueueDeclarePassive(name, durable, autoDelete, exclusive, noWait, args)
		return err
	})
	return queue, err
}

func (ch *ChannelWrapper) QueueInspect(name string) (queue Queue, err error) {
	err = ch.handle(func(channel *amqp.Channel) error {
		queue, err = channel.QueueInspect(name)
		return err
	})
	return queue, err
}

func (ch *ChannelWrapper) QueueBind(name string, key string, exchange string, noWait bool, args Table) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.QueueBind(name, key, exchange, noWait, args)
	})
}

func (ch *ChannelWrapper) QueueUnbind(name string, key string, exchange string, args Table) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.QueueUnbind(name, key, exchange, args)
	})
}

func (ch *ChannelWrapper) QueuePurge(name string, noWait bool) (count int, err error) {
	err = ch.handle(func(channel *amqp.Channel) error {
		count, err = channel.QueuePurge(name, noWait)
		return err
	})
	return count, err
}

func (ch *ChannelWrapper) QueueDelete(name string, ifUnused bool, ifEmpty bool, noWait bool) (count int, err error) {
	err = ch.handle(func(channel *amqp.Channel) error {
		count, err = channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
		return err
	})
	return count, err
}

func (ch *ChannelWrapper) Consume(queue string, consumer string, autoAck bool, exclusive bool,
	noLocal bool, noWait bool, args Table) (<-chan Delivery, error) {
	// TODO Add a wrapper to safely consume from new channel.
	return ch.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (ch *ChannelWrapper) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool,
	internal bool, noWait bool, args Table) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	})
}

func (ch *ChannelWrapper) ExchangeDeclarePassive(name string, kind string, durable bool, autoDelete bool,
	internal bool, noWait bool, args Table) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.ExchangeDeclarePassive(name, kind, durable, autoDelete, internal, noWait, args)
	})
}

func (ch *ChannelWrapper) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.ExchangeDelete(name, ifUnused, noWait)
	})
}

func (ch *ChannelWrapper) ExchangeBind(destination string, key string, source string, noWait bool, args Table) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.ExchangeBind(destination, key, source, noWait, args)
	})
}

func (ch *ChannelWrapper) ExchangeUnbind(destination string, key string, source string, noWait bool, args Table) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.ExchangeUnbind(destination, key, source, noWait, args)
	})
}

func (ch *ChannelWrapper) Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.Publish(exchange, key, mandatory, immediate, msg)
	})
}

func (ch *ChannelWrapper) Get(queue string, autoAck bool) (msg Delivery, ok bool, err error) {
	err = ch.handle(func(channel *amqp.Channel) error {
		msg, ok, err = channel.Get(queue, autoAck)
		return err
	})
	return msg, ok, err
}

func (ch *ChannelWrapper) Ack(tag uint64, multiple bool) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.Ack(tag, multiple)
	})
}

func (ch *ChannelWrapper) Nack(tag uint64, multiple bool, requeue bool) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.Nack(tag, multiple, requeue)
	})
}

func (ch *ChannelWrapper) Reject(tag uint64, requeue bool) error {
	return ch.handle(func(channel *amqp.Channel) error {
		return channel.Reject(tag, requeue)
	})
}

func (ch *ChannelWrapper) close(connection io.Closer) {
	err := connection.Close()
	if err != nil {
		ch.logger.Error("Failed to close connection")
		ch.observer.OnClose(err)
	}
}

func (ch *ChannelWrapper) handle(callback func(channel *amqp.Channel) error) error {
	var closed error

	handler := func() error {
		channel, err := ch.getChannel()
		if err == nil {
			// TODO (novln): Detect if we should release channel here.
			return callback(channel)
		}

		if errors.Cause(err) == ErrClientClosed {
			closed = err
			return nil
		}

		return err
	}

	err := ch.retrier.retry(handler)
	if closed != nil {
		return closed
	}

	return err
}

func (ch *ChannelWrapper) getChannel() (*amqp.Channel, error) {
	ch.mutex.RLock()
	channel := ch.channel
	ch.mutex.RUnlock()
	if channel != nil {
		return channel, nil
	}

	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	err := ch.newChannel(false)
	if err != nil {
		return nil, err
	}

	channel = ch.channel
	return channel, nil
}

func (ch *ChannelWrapper) releaseChannel() {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	if ch.channel != nil {
		ch.close(ch.channel)
	}
	ch.channel = nil
}

func (ch *ChannelWrapper) newChannel(lock bool) error {
	if lock {
		ch.mutex.Lock()
		defer ch.mutex.Unlock()
	}

	if ch.channel != nil {
		return nil
	}

	var closed error
	err := ch.retrier.retry(func() error {

		channel, err := ch.client.Channel()
		if err == nil {
			ch.channel = channel

			// Define Qos
			if ch.qos != nil {
				err = ch.channel.Qos(ch.qos.prefetchCount, ch.qos.prefetchSize, ch.qos.global)
				if err != nil {
					ch.close(ch.channel)
					ch.channel = nil
					return err
				}
			}

			return nil
		}

		if channel != nil {
			ch.close(channel)
		}

		if errors.Cause(err) == ErrClientClosed {
			closed = err
			return nil
		}

		return err
	})
	if err != nil {
		return err
	}
	if closed != nil {
		return closed
	}

	return nil
}

var _ Channel = (*ChannelWrapper)(nil)

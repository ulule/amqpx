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
	//
	// Qos controls how many messages or how many bytes the server will try to keep on
	// the network for consumers before receiving delivery acks. The intent of Qos is
	// to make sure the network buffers stay full between the server and client.
	//
	Qos(prefetchCount int, prefetchSize int, global bool) error
	//
	// Cancel stops deliveries to the consumer chan established in Channel.Consume and
	// identified by consumer.
	//
	Cancel(consumer string, noWait bool) error
	//
	// QueueDeclare declares a queue to hold messages and deliver to consumers.
	// Declaring creates a queue if it doesn't already exist, or ensures that an
	// existing queue matches the same parameters.
	//
	// Every queue declared gets a default binding to the empty exchange "" which has
	// the type "direct" with the routing key matching the queue's name.  With this
	// default binding, it is possible to publish messages that route directly to
	// this queue by publishing to "" with the routing key of the queue name.
	//
	QueueDeclare(name string, durable bool, autoDelete bool,
		exclusive bool, noWait bool, args Table) (Queue, error)
	//
	// QueueDeclarePassive is functionally and parametrically equivalent to
	// QueueDeclare, except that it sets the "passive" attribute to true. A passive
	// queue is assumed by RabbitMQ to already exist, and attempting to connect to a
	// non-existent queue will cause RabbitMQ to throw an exception. This function
	// can be used to test for the existence of a queue.
	//
	QueueDeclarePassive(name string, durable bool, autoDelete bool,
		exclusive bool, noWait bool, args Table) (Queue, error)
	//
	// QueueInspect passively declares a queue by name to inspect the current message
	// count and consumer count.
	//
	// Use this method to check how many messages ready for delivery reside in the queue,
	// how many consumers are receiving deliveries, and whether a queue by this
	// name already exists.
	//
	QueueInspect(name string) (Queue, error)
	//
	// QueueBind binds an exchange to a queue so that publishings to the exchange will
	// be routed to the queue when the publishing routing key matches the binding
	// routing key.
	//
	QueueBind(name string, key string, exchange string, noWait bool, args Table) error
	//
	// QueueUnbind removes a binding between an exchange and queue matching the key and
	// arguments.
	//
	// It is possible to send and empty string for the exchange name which means to
	// unbind the queue from the default exchange.
	//
	QueueUnbind(name string, key string, exchange string, args Table) error
	//
	// QueuePurge removes all messages from the named queue which are not waiting to
	// be acknowledged.  Messages that have been delivered but have not yet been
	// acknowledged will not be removed.
	//
	// When successful, returns the number of messages purged.
	//
	QueuePurge(name string, noWait bool) (int, error)
	//
	// QueueDelete removes the queue from the server including all bindings then
	// purges the messages based on server configuration, returning the number of
	// messages purged.
	//
	QueueDelete(name string, ifUnused bool, ifEmpty bool, noWait bool) (int, error)
	//
	// Consume immediately starts delivering queued messages.
	//
	// Begin receiving on the returned chan Delivery before any other operation on the
	// Connection or Channel.
	//
	Consume(queue string, consumer string, autoAck bool, exclusive bool,
		noLocal bool, noWait bool, args Table) (<-chan Delivery, error)
	//
	// ExchangeDeclare declares an exchange on the server. If the exchange does not
	// already exist, the server will create it.  If the exchange exists, the server
	// verifies that it is of the provided type, durability and auto-delete flags.
	//
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool,
		internal bool, noWait bool, args Table) error
	//
	// ExchangeDeclarePassive is functionally and parametrically equivalent to
	// ExchangeDeclare, except that it sets the "passive" attribute to true. A passive
	// exchange is assumed by RabbitMQ to already exist, and attempting to connect to a
	// non-existent exchange will cause RabbitMQ to throw an exception. This function
	// can be used to detect the existence of an exchange.
	//
	ExchangeDeclarePassive(name string, kind string, durable bool, autoDelete bool,
		internal bool, noWait bool, args Table) error
	//
	// ExchangeDelete removes the named exchange from the server. When an exchange is
	// deleted all queue bindings on the exchange are also deleted.  If this exchange
	// does not exist, the channel will be closed with an error.
	//
	ExchangeDelete(name string, ifUnused, noWait bool) error
	//
	// ExchangeBind binds an exchange to another exchange to create inter-exchange
	// routing topologies on the server.  This can decouple the private topology and
	// routing exchanges from exchanges intended solely for publishing endpoints.
	//
	// Binding two exchanges with identical arguments will not create duplicate
	// bindings.
	//
	ExchangeBind(destination string, key string, source string, noWait bool, args Table) error
	//
	// ExchangeUnbind unbinds the destination exchange from the source exchange on the
	// server by removing the routing key between them.  This is the inverse of
	// ExchangeBind.  If the binding does not currently exist, an error will be
	// returned.
	//
	ExchangeUnbind(destination string, key string, source string, noWait bool, args Table) error
	//
	// Publish sends a Publishing from the client to an exchange on the server.
	//
	// When you want a single message to be delivered to a single queue, you can
	// publish to the default exchange with the routingKey of the queue name.  This is
	// because every declared queue gets an implicit route to the default exchange.
	//
	Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error
	//
	// Get synchronously receives a single Delivery from the head of a queue from the
	// server to the client.  In almost all cases, using Channel.Consume will be
	// preferred.
	//
	// If there was a delivery waiting on the queue and that delivery was received, the
	// second return value will be true.  If there was no delivery waiting or an error
	// occurred, the ok bool will be false.
	//
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
			thr := callback(channel)
			if thr != nil {
				ch.releaseChannel()
			}
			return thr
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

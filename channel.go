package amqpx

import (
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Channel describes an amqp channel.
type Channel interface {
	Close() error
	NotifyClose(c chan *Error) chan *Error
	NotifyFlow(c chan bool) chan bool
	NotifyReturn(c chan Return) chan Return
	NotifyCancel(c chan string) chan string
	NotifyConfirm(ack chan uint64, nack chan uint64) (chan uint64, chan uint64)
	NotifyPublish(confirm chan Confirmation) chan Confirmation
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
	Tx() error
	TxCommit() error
	TxRollback() error
	Flow(active bool) error
	Confirm(noWait bool) error
	Recover(requeue bool) error
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple bool, requeue bool) error
	Reject(tag uint64, requeue bool) error
}

// ChannelWrapper wraps a amqp channel to provide a retry mechanism.
type ChannelWrapper struct {
	mutex    sync.RWMutex
	client   Client
	observer Observer
	logger   Logger
	retrier  retrier
	channel  *amqp.Channel
}

func NewChannelWrapper(client Client, retrier retrier) Channel {
	return &ChannelWrapper{
		client:   client,
		logger:   client.getLogger(),
		observer: client.getObserver(),
		retrier:  retrier,
	}
}

// TODO (novln): Decide what and how we obtain a channel on this wrapper.

func (ch *ChannelWrapper) Close() error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	if ch.channel != nil {
		err := ch.channel.Close()
		if err != nil {
			return errors.Wrap(err, ErrMessageCannotCloseChannel)
		}
		ch.channel = nil
	}

	return nil
}

func (ch *ChannelWrapper) NotifyClose(c chan *Error) chan *Error {
	return ch.channel.NotifyClose(c)
}

func (ch *ChannelWrapper) NotifyFlow(c chan bool) chan bool {
	return ch.channel.NotifyFlow(c)
}

func (ch *ChannelWrapper) NotifyReturn(c chan Return) chan Return {
	return ch.channel.NotifyReturn(c)
}

func (ch *ChannelWrapper) NotifyCancel(c chan string) chan string {
	return ch.channel.NotifyCancel(c)
}

func (ch *ChannelWrapper) NotifyConfirm(ack chan uint64, nack chan uint64) (chan uint64, chan uint64) {
	return ch.channel.NotifyConfirm(ack, nack)
}

func (ch *ChannelWrapper) NotifyPublish(confirm chan Confirmation) chan Confirmation {
	return ch.channel.NotifyPublish(confirm)
}

func (ch *ChannelWrapper) Qos(prefetchCount int, prefetchSize int, global bool) error {
	return ch.channel.Qos(prefetchCount, prefetchSize, global)
}

func (ch *ChannelWrapper) Cancel(consumer string, noWait bool) error {
	return ch.channel.Cancel(consumer, noWait)
}

func (ch *ChannelWrapper) QueueDeclare(name string, durable bool, autoDelete bool,
	exclusive bool, noWait bool, args Table) (Queue, error) {
	return ch.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (ch *ChannelWrapper) QueueDeclarePassive(name string, durable bool, autoDelete bool,
	exclusive bool, noWait bool, args Table) (Queue, error) {
	return ch.channel.QueueDeclarePassive(name, durable, autoDelete, exclusive, noWait, args)
}

func (ch *ChannelWrapper) QueueInspect(name string) (Queue, error) {
	return ch.channel.QueueInspect(name)
}

func (ch *ChannelWrapper) QueueBind(name string, key string, exchange string, noWait bool, args Table) error {
	return ch.channel.QueueBind(name, key, exchange, noWait, args)
}

func (ch *ChannelWrapper) QueueUnbind(name string, key string, exchange string, args Table) error {
	return ch.channel.QueueUnbind(name, key, exchange, args)
}

func (ch *ChannelWrapper) QueuePurge(name string, noWait bool) (int, error) {
	return ch.channel.QueuePurge(name, noWait)
}

func (ch *ChannelWrapper) QueueDelete(name string, ifUnused bool, ifEmpty bool, noWait bool) (int, error) {
	return ch.channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
}

func (ch *ChannelWrapper) Consume(queue string, consumer string, autoAck bool, exclusive bool,
	noLocal bool, noWait bool, args Table) (<-chan Delivery, error) {
	return ch.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (ch *ChannelWrapper) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool,
	internal bool, noWait bool, args Table) error {
	return ch.channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (ch *ChannelWrapper) ExchangeDeclarePassive(name string, kind string, durable bool, autoDelete bool,
	internal bool, noWait bool, args Table) error {
	return ch.channel.ExchangeDeclarePassive(name, kind, durable, autoDelete, internal, noWait, args)
}

func (ch *ChannelWrapper) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return ch.channel.ExchangeDelete(name, ifUnused, noWait)
}

func (ch *ChannelWrapper) ExchangeBind(destination string, key string, source string, noWait bool, args Table) error {
	return ch.channel.ExchangeUnbind(destination, key, source, noWait, args)
}

func (ch *ChannelWrapper) ExchangeUnbind(destination string, key string, source string, noWait bool, args Table) error {
	return ch.channel.ExchangeUnbind(destination, key, source, noWait, args)
}

func (ch *ChannelWrapper) Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error {
	handler := func() error {
		return ch.channel.Publish(exchange, key, mandatory, immediate, msg)
	}
	return ch.retrier.retry(handler)
}

func (ch *ChannelWrapper) Get(queue string, autoAck bool) (msg Delivery, ok bool, err error) {
	return ch.channel.Get(queue, autoAck)
}

func (ch *ChannelWrapper) Tx() error {
	return ch.channel.Tx()
}

func (ch *ChannelWrapper) TxCommit() error {
	return ch.channel.TxCommit()
}

func (ch *ChannelWrapper) TxRollback() error {
	return ch.channel.TxRollback()
}

func (ch *ChannelWrapper) Flow(active bool) error {
	return ch.channel.Flow(active)
}

func (ch *ChannelWrapper) Confirm(noWait bool) error {
	return ch.channel.Confirm(noWait)
}

func (ch *ChannelWrapper) Recover(requeue bool) error {
	return ch.channel.Recover(requeue)
}

func (ch *ChannelWrapper) Ack(tag uint64, multiple bool) error {
	return ch.channel.Ack(tag, multiple)
}

func (ch *ChannelWrapper) Nack(tag uint64, multiple bool, requeue bool) error {
	return ch.channel.Nack(tag, multiple, requeue)
}

func (ch *ChannelWrapper) Reject(tag uint64, requeue bool) error {
	return ch.channel.Reject(tag, requeue)
}

func (ch *ChannelWrapper) releaseChannel() {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	if ch.channel != nil {
		ch.close(ch.channel)
	}
	ch.channel = nil
}

func (ch *ChannelWrapper) close(connection io.Closer) {
	err := connection.Close()
	if err != nil {
		ch.logger.Error("Failed to close connection")
		ch.observer.OnClose(err)
	}
}

func (ch *ChannelWrapper) newChannel() error {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()

	if ch.channel != nil {
		return nil
	}

	var closed error
	err := ch.retrier.retry(func() error {

		channel, err := ch.client.newChannel()
		if err == nil {
			ch.channel = channel
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

package amqpx

import (
	"github.com/streadway/amqp"
	"github.com/ulule/kyu/try"
)

type channelRetry struct {
	channel *amqp.Channel
	retryOptions
}

func (ch channelRetry) Close() error {
	return ch.channel.Close()
}

func (ch channelRetry) NotifyClose(c chan *Error) chan *Error {
	return ch.channel.NotifyClose(c)
}

func (ch channelRetry) NotifyFlow(c chan bool) chan bool {
	return ch.channel.NotifyFlow(c)
}

func (ch channelRetry) NotifyReturn(c chan Return) chan Return {
	return ch.channel.NotifyReturn(c)
}

func (ch channelRetry) NotifyCancel(c chan string) chan string {
	return ch.channel.NotifyCancel(c)
}

func (ch channelRetry) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	return ch.channel.NotifyConfirm(ack, nack)
}

func (ch channelRetry) NotifyPublish(confirm chan Confirmation) chan Confirmation {
	return ch.channel.NotifyPublish(confirm)
}

func (ch channelRetry) Qos(prefetchCount, prefetchSize int, global bool) error {
	return ch.channel.Qos(prefetchCount, prefetchSize, global)
}

func (ch channelRetry) Cancel(consumer string, noWait bool) error {
	return ch.channel.Cancel(consumer, noWait)
}

func (ch channelRetry) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error) {
	return ch.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (ch channelRetry) QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error) {
	return ch.channel.QueueDeclarePassive(name, durable, autoDelete, exclusive, noWait, args)
}

func (ch channelRetry) QueueInspect(name string) (Queue, error) {
	return ch.channel.QueueInspect(name)
}

func (ch channelRetry) QueueBind(name, key, exchange string, noWait bool, args Table) error {
	return ch.channel.QueueBind(name, key, exchange, noWait, args)
}

func (ch channelRetry) QueueUnbind(name, key, exchange string, args Table) error {
	return ch.channel.QueueUnbind(name, key, exchange, args)
}

func (ch channelRetry) QueuePurge(name string, noWait bool) (int, error) {
	return ch.channel.QueuePurge(name, noWait)
}

func (ch channelRetry) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	return ch.channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
}

func (ch channelRetry) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error) {
	return ch.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

func (ch channelRetry) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error {
	return ch.channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (ch channelRetry) ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args Table) error {
	return ch.channel.ExchangeDeclarePassive(name, kind, durable, autoDelete, internal, noWait, args)
}

func (ch channelRetry) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return ch.channel.ExchangeDelete(name, ifUnused, noWait)
}

func (ch channelRetry) ExchangeBind(destination, key, source string, noWait bool, args Table) error {
	return ch.channel.ExchangeUnbind(destination, key, source, noWait, args)
}

func (ch channelRetry) ExchangeUnbind(destination, key, source string, noWait bool, args Table) error {
	return ch.channel.ExchangeUnbind(destination, key, source, noWait, args)
}

// Publish overrides amqp.Channel.Publish method to provide a retry mechanism.
func (ch *channelRetry) Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error {
	return try.Try(func() error {
		return ch.channel.Publish(exchange, key, mandatory, immediate, msg)
	}, try.Option{
		InitialInterval: ch.retryInitialInterval,
		MaxInterval:     ch.retryMaxInterval,
		MaxElapsedTime:  ch.retryMaxElapsedTime,
	})
}

func (ch channelRetry) Get(queue string, autoAck bool) (msg Delivery, ok bool, err error) {
	return ch.channel.Get(queue, autoAck)
}

func (ch channelRetry) Tx() error {
	return ch.channel.Tx()
}

func (ch channelRetry) TxCommit() error {
	return ch.channel.TxCommit()
}

func (ch channelRetry) TxRollback() error {
	return ch.channel.TxRollback()
}

func (ch channelRetry) Flow(active bool) error {
	return ch.channel.Flow(active)
}

func (ch channelRetry) Confirm(noWait bool) error {
	return ch.channel.Confirm(noWait)
}

func (ch channelRetry) Recover(requeue bool) error {
	return ch.channel.Recover(requeue)
}

func (ch channelRetry) Ack(tag uint64, multiple bool) error {
	return ch.channel.Ack(tag, multiple)
}

func (ch channelRetry) Nack(tag uint64, multiple bool, requeue bool) error {
	return ch.channel.Nack(tag, multiple, requeue)
}

func (ch channelRetry) Reject(tag uint64, requeue bool) error {
	return ch.channel.Reject(tag, requeue)
}

var _ Channel = (*channelRetry)(nil)

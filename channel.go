package amqpx

import "github.com/streadway/amqp"

// Channel represents an amqp channel.
type Channel interface {
	Close() error
	NotifyClose(c chan *Error) chan *Error
	NotifyFlow(c chan bool) chan bool
	NotifyReturn(c chan Return) chan Return
	NotifyCancel(c chan string) chan string
	NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64)
	NotifyPublish(confirm chan Confirmation) chan Confirmation
	Qos(prefetchCount, prefetchSize int, global bool) error
	Cancel(consumer string, noWait bool) error
	QueueDeclare(
		name string,
		durable,
		autoDelete,
		exclusive,
		noWait bool,
		args Table) (Queue, error)
	QueueDeclarePassive(
		name string,
		durable,
		autoDelete,
		exclusive,
		noWait bool,
		args Table) (Queue, error)
	QueueInspect(name string) (Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args Table) error
	QueueUnbind(name, key, exchange string, args Table) error
	QueuePurge(name string, noWait bool) (int, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
	Consume(
		queue,
		consumer string,
		autoAck,
		exclusive,
		noLocal,
		noWait bool,
		args Table) (<-chan Delivery, error)
	ExchangeDeclare(
		name,
		kind string,
		durable,
		autoDelete,
		internal,
		noWait bool,
		args Table) error
	ExchangeDeclarePassive(
		name,
		kind string,
		durable,
		autoDelete,
		internal,
		noWait bool,
		args Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeBind(
		destination,
		key,
		source string,
		noWait bool,
		args Table) error
	ExchangeUnbind(
		destination,
		key,
		source string,
		noWait bool,
		args Table) error
	Publish(
		exchange string,
		key string,
		mandatory bool,
		immediate bool,
		msg Publishing) error
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

type channelWrapper struct {
	channel *amqp.Channel
	retrier retrier
}

func newChannel(ch *amqp.Channel, rt retrier) Channel {
	return &channelWrapper{
		channel: ch,
		retrier: rt,
	}
}

func (ch channelWrapper) Close() error {
	return ch.channel.Close()
}

func (ch channelWrapper) NotifyClose(c chan *Error) chan *Error {
	return ch.channel.NotifyClose(c)
}

func (ch channelWrapper) NotifyFlow(c chan bool) chan bool {
	return ch.channel.NotifyFlow(c)
}

func (ch channelWrapper) NotifyReturn(c chan Return) chan Return {
	return ch.channel.NotifyReturn(c)
}

func (ch channelWrapper) NotifyCancel(c chan string) chan string {
	return ch.channel.NotifyCancel(c)
}

func (ch channelWrapper) NotifyConfirm(
	ack,
	nack chan uint64) (chan uint64, chan uint64) {

	return ch.channel.NotifyConfirm(ack, nack)
}

func (ch channelWrapper) NotifyPublish(
	confirm chan Confirmation) chan Confirmation {

	return ch.channel.NotifyPublish(confirm)
}

func (ch channelWrapper) Qos(
	prefetchCount,
	prefetchSize int,
	global bool) error {

	return ch.channel.Qos(prefetchCount, prefetchSize, global)
}

func (ch channelWrapper) Cancel(consumer string, noWait bool) error {
	return ch.channel.Cancel(consumer, noWait)
}

func (ch channelWrapper) QueueDeclare(
	name string,
	durable,
	autoDelete,
	exclusive,
	noWait bool,
	args Table) (Queue, error) {

	return ch.channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args)
}

func (ch channelWrapper) QueueDeclarePassive(
	name string,
	durable,
	autoDelete,
	exclusive,
	noWait bool,
	args Table) (Queue, error) {

	return ch.channel.QueueDeclarePassive(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args)
}

func (ch channelWrapper) QueueInspect(name string) (Queue, error) {
	return ch.channel.QueueInspect(name)
}

func (ch channelWrapper) QueueBind(
	name,
	key,
	exchange string,
	noWait bool,
	args Table) error {

	return ch.channel.QueueBind(name, key, exchange, noWait, args)
}

func (ch channelWrapper) QueueUnbind(
	name,
	key,
	exchange string, args Table) error {

	return ch.channel.QueueUnbind(name, key, exchange, args)
}

func (ch channelWrapper) QueuePurge(name string, noWait bool) (int, error) {
	return ch.channel.QueuePurge(name, noWait)
}

func (ch channelWrapper) QueueDelete(
	name string,
	ifUnused,
	ifEmpty,
	noWait bool) (int, error) {

	return ch.channel.QueueDelete(name, ifUnused, ifEmpty, noWait)
}

func (ch channelWrapper) Consume(
	queue,
	consumer string,
	autoAck,
	exclusive,
	noLocal,
	noWait bool,
	args Table) (<-chan Delivery, error) {

	return ch.channel.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args)
}

func (ch channelWrapper) ExchangeDeclare(
	name,
	kind string,
	durable,
	autoDelete,
	internal,
	noWait bool,
	args Table) error {

	return ch.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args)
}

func (ch channelWrapper) ExchangeDeclarePassive(
	name,
	kind string,
	durable,
	autoDelete,
	internal,
	noWait bool,
	args Table) error {

	return ch.channel.ExchangeDeclarePassive(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args)
}

func (ch channelWrapper) ExchangeDelete(
	name string,
	ifUnused,
	noWait bool) error {

	return ch.channel.ExchangeDelete(name, ifUnused, noWait)
}

func (ch channelWrapper) ExchangeBind(
	destination,
	key,
	source string,
	noWait bool,
	args Table) error {

	return ch.channel.ExchangeUnbind(destination, key, source, noWait, args)
}

func (ch channelWrapper) ExchangeUnbind(
	destination,
	key,
	source string,
	noWait bool,
	args Table) error {

	return ch.channel.ExchangeUnbind(destination, key, source, noWait, args)
}

// Publish overrides amqp.Channel.Publish method to provide a retry mechanism.
func (ch *channelWrapper) Publish(
	exchange string,
	key string,
	mandatory bool,
	immediate bool,
	msg Publishing) error {

	return ch.retrier.retry(func() error {
		return ch.channel.Publish(
			exchange,
			key,
			mandatory,
			immediate,
			msg)
	})
}

func (ch channelWrapper) Get(
	queue string,
	autoAck bool) (msg Delivery, ok bool, err error) {

	return ch.channel.Get(queue, autoAck)
}

func (ch channelWrapper) Tx() error {
	return ch.channel.Tx()
}

func (ch channelWrapper) TxCommit() error {
	return ch.channel.TxCommit()
}

func (ch channelWrapper) TxRollback() error {
	return ch.channel.TxRollback()
}

func (ch channelWrapper) Flow(active bool) error {
	return ch.channel.Flow(active)
}

func (ch channelWrapper) Confirm(noWait bool) error {
	return ch.channel.Confirm(noWait)
}

func (ch channelWrapper) Recover(requeue bool) error {
	return ch.channel.Recover(requeue)
}

func (ch channelWrapper) Ack(tag uint64, multiple bool) error {
	return ch.channel.Ack(tag, multiple)
}

func (ch channelWrapper) Nack(tag uint64, multiple bool, requeue bool) error {
	return ch.channel.Nack(tag, multiple, requeue)
}

func (ch channelWrapper) Reject(tag uint64, requeue bool) error {
	return ch.channel.Reject(tag, requeue)
}

var _ Channel = (*channelWrapper)(nil)

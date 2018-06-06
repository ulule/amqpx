package amqpx

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

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
	*amqp.Channel
	retrier retrier
}

func newChannel(ch *amqp.Channel, rt retrier) Channel {
	return &channelWrapper{
		Channel: ch,
		retrier: rt,
	}
}

// Publish overrides amqp.Channel.Publish method to provide a retry mechanism.
func (ch *channelWrapper) Publish(
	exchange string,
	key string,
	mandatory bool,
	immediate bool,
	msg Publishing) error {

	return ch.retrier.retry(func() error {
		return ch.Channel.Publish(
			exchange,
			key,
			mandatory,
			immediate,
			msg)
	})
}

func openChannel(conn *amqp.Connection, retry retrier, obs Observer) (Channel, error) {
	var (
		err error
		ch  *amqp.Channel
	)

	err = retry.retry(func() error {
		ch, err = conn.Channel()
		if err != nil && err != amqp.ErrClosed {
			if ch != nil {
				err = ch.Close()
				if err != nil {
					obs.OnClose(err)
				}
			}
			return errors.Wrap(err, ErrOpenChannel.Error())
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "exceeded retries to open connection channel")
	}

	return newChannel(ch, retry), nil
}

var _ Channel = (*channelWrapper)(nil)

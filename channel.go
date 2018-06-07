package amqpx

import (
	"fmt"

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

// channelWrapper wraps a amqp channel to provide a retry mechanism.
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
func (ch *channelWrapper) Publish(exchange string, key string, mandatory bool, immediate bool, msg Publishing) error {
	handler := func() error {
		return ch.Channel.Publish(exchange, key, mandatory, immediate, msg)
	}
	return ch.retrier.retry(handler)
}

func openChannel(conn *amqp.Connection, retryOpts retrierOptions, obs Observer, logger Logger) (Channel, error) {
	var (
		err   error
		ch    *amqp.Channel
		retry = newRetrier(retryOpts)
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

			logger.Error(
				fmt.Sprintf("Retry to obtain a channel for connection %s",
					conn.LocalAddr()))

			return errors.Wrap(err, ErrMessageCannotOpenChannel)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, ErrMessageRetryExceeded)
	}

	logger.Debug(fmt.Sprintf("Opened channel on connection %s", conn.LocalAddr()))

	return newChannel(ch, retry), nil
}

var _ Channel = (*channelWrapper)(nil)

package amqpx

import (
	"github.com/streadway/amqp"
	"github.com/ulule/kyu/try"
)

// Channel is a wrapper around amqp.Channel to provide retry feature.
type Channel struct {
	*amqp.Channel
	retryOptions
}

// Publish overrides amqp.Channel.Publish method to provide a retry mechanism.
func (ch *Channel) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return try.Try(func() error {
		return ch.Channel.Publish(exchange, key, mandatory, immediate, msg)
	}, try.Option{
		InitialInterval: ch.retryInitialInterval,
		MaxInterval:     ch.retryMaxInterval,
		MaxElapsedTime:  ch.retryMaxElapsedTime,
	})
}

package amqpx_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/ulule/amqpx"
)

const (
	warpExchange = "amq.topic"
)

func NewMessage(body []byte) amqp.Publishing {
	return amqp.Publishing{
		DeliveryMode:    amqp.Persistent,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Body:            body,
	}
}

type Event struct {
	Message string
}

type Emitter struct {
	client amqpx.Client
}

func NewEmitter(client amqpx.Client) *Emitter {
	return &Emitter{client: client}
}

func (emitter *Emitter) setDirectChannel(channel *amqp.Channel, topic string) error {
	directQueue, err := channel.QueueDeclare(
		topic,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot declare topic")
	}

	err = channel.QueueBind(
		directQueue.Name, // queue name
		directQueue.Name, // routing key
		warpExchange,     // exchange
		false,            // noWait
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot bind exchange")
	}

	return nil
}

// Publish emits a message to a given topic
func (emitter *Emitter) Publish(topic string, payload amqp.Publishing) error {
	channel, err := emitter.client.Channel()
	if channel != nil {
		defer func() {
			thr := channel.Close()
			_ = thr
		}()
	}
	if err != nil {
		return errors.Wrap(err, "cannot acquire a new channel from client")
	}

	err = emitter.setDirectChannel(channel, topic)
	if err != nil {
		return errors.Wrapf(err, "cannot acquire a channel for topic: %s", topic)
	}

	err = channel.Publish(
		"",    // default exchange
		topic, // queue name
		false, // mandatory
		false, // immediate
		payload,
	)
	if err != nil {
		return errors.Wrapf(err, "cannot publish on queue: %s", topic)
	}

	return nil
}

func (emitter *Emitter) Close() error {
	return emitter.client.Close()
}

type Handler func(message []byte)

type Receiver struct {
	client   amqpx.Client
	name     string
	consumer string
}

func NewReceiver(client amqpx.Client, name string) *Receiver {
	return &Receiver{
		client:   client,
		name:     name,
		consumer: fmt.Sprintf("%s_%d", name, rand.Intn(100)),
	}
}

func (receiver *Receiver) Start(handler Handler) error {
	channel, err := receiver.client.Channel()
	if channel != nil {
		defer func() {
			thr := channel.Close()
			_ = thr
		}()
	}
	if err != nil {
		return errors.Wrap(err, "cannot acquire a channel")
	}

	_, err = channel.QueueDeclare(
		receiver.name,
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return errors.Wrapf(err, "cannot declare queue: %s", receiver.name)
	}

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return errors.Wrap(err, "cannot configure channel")
	}

	deliveries, err := channel.Consume(
		receiver.name,     // queue name
		receiver.consumer, // unique consumer id
		false,             // noAck
		false,             // exclusive
		false,             // noLocal
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		return errors.Wrap(err, "cannot create a message receiver")
	}

	for {
		select {
		case delivery := <-deliveries:
			handler(delivery.Body)
			err := delivery.Ack(false)
			if err != nil {
				return err
			}
		case <-time.After(30 * time.Second):
			return nil
		}
	}
}

type Producer struct {
	runner   *Runner
	messages []string
	buffer   chan string
	wg       *sync.WaitGroup
}

func (producer *Producer) Start() {
	producer.buffer = make(chan string, 256)
	producer.wg = &sync.WaitGroup{}

	go func() {
		for _, message := range producer.messages {
			producer.buffer <- message
		}
		close(producer.buffer)
	}()
}

func (producer *Producer) NewEmitter(client amqpx.Client, topic string) {
	producer.wg.Add(1)
	go func() {
		emitter := NewEmitter(client)
		defer producer.wg.Done()
		for {

			message, ok := <-producer.buffer
			if !ok {
				return
			}

			event := &Event{Message: message}
			json, err := json.Marshal(event)
			producer.runner.NoError(err)

			err = emitter.Publish(topic, NewMessage(json))
			producer.runner.NoError(err)
		}
	}()
}

func (producer *Producer) Wait() {
	producer.wg.Wait()
}

type Consumer struct {
	runner   *Runner
	messages []string
	buffer   chan string
	wg       *sync.WaitGroup
	killed   chan struct{}
}

func (consumer *Consumer) Start() {
	consumer.buffer = make(chan string, 256)
	consumer.wg = &sync.WaitGroup{}
	consumer.killed = make(chan struct{})
	go func() {
		for {
			message, ok := <-consumer.buffer
			if !ok {
				consumer.killed <- struct{}{}
				return
			}
			consumer.messages = append(consumer.messages, message)
		}
	}()
}

func (consumer *Consumer) NewReceiver(client amqpx.Client, topic string) {
	consumer.wg.Add(1)
	go func() {
		receiver := NewReceiver(client, topic)
		defer consumer.wg.Done()
		err := receiver.Start(func(body []byte) {
			event := &Event{}
			err := json.Unmarshal(body, event)
			consumer.runner.NoError(err)
			consumer.buffer <- event.Message
		})
		consumer.runner.NoError(err)
	}()
}

func (consumer *Consumer) Wait() {
	consumer.wg.Wait()
	close(consumer.buffer)
	<-consumer.killed
}

func GenerateMessages() []string {
	buffer := []string{}
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")
	length := len(runes)
	lines := 3000 + rand.Intn(2000)

	for i := 0; i < lines; i++ {
		line := &bytes.Buffer{}
		z := 32 + rand.Intn(255)
		for y := 0; y < z; y++ {
			line.WriteRune(runes[rand.Intn(length)])
		}
		buffer = append(buffer, line.String())
	}

	return buffer
}

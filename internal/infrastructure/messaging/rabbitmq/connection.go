package rmq_rpc

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Params struct {
	URL      string
	WaitTime time.Duration
	Attempts int
}

type Connection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Deliveries map[string]<-chan amqp.Delivery

	ConsumerExchange string
	Topics           []string
	Params
}

func NewConnection(consumerExchange string, topics []string, params Params) *Connection {
	return &Connection{
		ConsumerExchange: consumerExchange,
		Topics:           topics,
		Params:           params,
	}
}

func (c *Connection) AttemptConnect(connector func() error) error {
	var err error

	for i := c.Attempts; i > 0; i-- {
		if err = connector(); err == nil {
			break
		}

		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("rmq_rpc - AttemptConnect - c.connect: %w", err)
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("c.Connection.Channel: %w", err)
	}

	return nil
}

func (c *Connection) ConnectWriter() error {
	var err error

	err = c.connect()
	if err != nil {
		return err
	}

	for _, topic := range c.Topics {
		err = c.Channel.ExchangeDeclare(
			topic,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("c.Channel.ExchangeDeclare for topic %s: %w", topic, err)
		}
	}

	return nil
}

func (c *Connection) ConnectReader() error {
	var err error

	err = c.connect()
	if err != nil {
		return err
	}

	deliveries := make(map[string]<-chan amqp.Delivery)

	for _, topic := range c.Topics {
		err = c.Channel.ExchangeDeclare(
			topic,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("c.Channel.ExchangeDeclare for topic %s: %w", topic, err)
		}

		queue, err := c.Channel.QueueDeclare(
			"",
			true,
			false,
			true,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("c.Channel.QueueDeclare for topic %s: %w", topic, err)
		}

		err = c.Channel.QueueBind(
			queue.Name,
			"",
			topic,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("c.Channel.QueueBind for topic %s: %w", topic, err)
		}

		delivery, err := c.Channel.Consume(
			queue.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("c.Channel.Consume for topic %s: %w", topic, err)
		}

		deliveries[topic] = delivery
	}

	c.Deliveries = deliveries
	return nil
}

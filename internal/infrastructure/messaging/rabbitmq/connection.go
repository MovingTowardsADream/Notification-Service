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
	Delivery   <-chan amqp.Delivery

	ConsumerExchange string
	Params
}

func NewConnection(consumerExchange string, params Params) *Connection {
	conn := &Connection{
		ConsumerExchange: consumerExchange,
		Params:           params,
	}

	return conn
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

	err = c.Channel.ExchangeDeclare(
		c.ConsumerExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.ExchangeDeclare: %w", err)
	}

	return nil
}

func (c *Connection) ConnectReader() error {
	var err error

	err = c.connect()
	if err != nil {
		return err
	}

	err = c.Channel.ExchangeDeclare(
		c.ConsumerExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.ExchangeDeclare: %w", err)
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
		return fmt.Errorf("c.Channel.QueueDeclare: %w", err)
	}

	err = c.Channel.QueueBind(
		queue.Name,
		"",
		c.ConsumerExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.QueueBind: %w", err)
	}

	c.Delivery, err = c.Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.Consume: %w", err)
	}

	return nil
}

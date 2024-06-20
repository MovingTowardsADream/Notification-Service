package rmq_rpc

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Config struct {
	URL      string
	WaitTime time.Duration
	Attempts int
}

type Connection struct {
	ConsumerExchange string
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Delivery   <-chan amqp.Delivery // Канал для получения доставленных сообщений
}

func NewConnectionRabbitMQ(consumerExchange string, cfg Config) *Connection {
	conn := &Connection{
		ConsumerExchange: consumerExchange,
		Config:           cfg,
	}

	return conn
}

func (c *Connection) AttemptConnect() error {
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
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

	// Establishing a connection to the RabbitMQ server
	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	// Open a communication channel
	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("c.Connection.Channel: %w", err)
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
		return fmt.Errorf("c.Connection.Channel: %w", err)
	}

	// Declares a queue to receive messages from the exchange
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

	// Связывает очередь с fanout exchange, используя пустой routing key
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

	// Starts consuming messages from the queue, providing a c.Delivery channel
	// through which received messages will be available.
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

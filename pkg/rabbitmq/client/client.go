package client

import (
	rmq_rpc "Notification_Service/pkg/rabbitmq"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"time"
)

var ErrConnectionClosed = errors.New("rmq_rpc client - Client - RemoteCall - Connection closed")

const (
	_defaultWaitTime = 2 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second

	Success = "success"
)

// RPC-клиента
type Client struct {
	conn           *rmq_rpc.Connection
	serverExchange string
	error          chan error
	stop           chan struct{}

	timeout time.Duration
}

// Creating a new RPC client
func NewRabbitMQClient(url, serverExchange, clientExchange string, opts ...Option) (*Client, error) {
	cfg := rmq_rpc.Config{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	c := &Client{
		conn:           rmq_rpc.NewConnectionRabbitMQ(clientExchange, cfg),
		serverExchange: serverExchange,
		error:          make(chan error),
		stop:           make(chan struct{}),
		timeout:        _defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc client - NewClient - c.conn.AttemptConnect: %w", err)
	}

	return c, nil
}

// Publish message in RabbitMQ
func (c *Client) publish(corrID, handler string, request interface{}) error {
	var (
		requestBody []byte
		err         error
	)

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return fmt.Errorf("publish - json.Marshal: %w", err)
		}
	}

	err = c.conn.Channel.Publish(c.serverExchange, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.conn.ConsumerExchange,
			Type:          handler,
			Body:          requestBody,
		})
	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}

func (c *Client) RemoteCall(ctx context.Context, handler string, request interface{}) error {
	// Если клиент не завершился
	select {
	case <-c.stop:
		time.Sleep(c.timeout)
		select {
		case <-c.stop:
			return ErrConnectionClosed
		default:
		}
	default:
	}

	corrID := uuid.New().String()
	err := c.publish(corrID, handler, request)
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - RemoteCall - c.publish: %w", err)
	}
	return nil
}

func (c *Client) Notify() <-chan error {
	return c.error
}

func (c *Client) Shutdown() error {
	select {
	case <-c.error:
		return nil
	default:
	}

	close(c.stop)
	time.Sleep(c.timeout)

	err := c.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("rmq_rpc client - Client - Shutdown - c.Connection.Close: %w", err)
	}

	return nil
}

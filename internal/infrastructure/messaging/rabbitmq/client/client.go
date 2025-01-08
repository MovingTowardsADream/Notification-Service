package rmq_client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"

	rmq_rpc "Notification_Service/internal/infrastructure/messaging/rabbitmq"
)

const (
	_defaultWaitTime = 2 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

type Client struct {
	conn *rmq_rpc.Connection

	serverExchange string
	timeout        time.Duration

	error chan error
	stop  chan struct{}
}

func New(url, serverExchange, clientExchange string, topics []string, opts ...Option) (*Client, error) {
	cfg := rmq_rpc.Params{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	c := &Client{
		conn:           rmq_rpc.NewConnection(clientExchange, topics, cfg),
		serverExchange: serverExchange,
		error:          make(chan error),
		stop:           make(chan struct{}),
		timeout:        _defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect(c.conn.ConnectWriter)
	if err != nil {
		return nil, fmt.Errorf("rmq_rpc client - NewClient - c.conn.AttemptConnect: %w", err)
	}

	return c, nil
}

func (c *Client) publish(corrID, handler, topic string, request any) error {
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

	err = c.conn.Channel.Publish(topic, "", false, false,
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

func (c *Client) RemoteCall(ctx context.Context, handler, topic string, request any) error {
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
	err := c.publish(corrID, handler, topic, request)
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

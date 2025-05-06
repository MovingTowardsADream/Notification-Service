package rmqclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/domain/models"
	rmqrpc "Notification_Service/internal/infrastructure/messaging/rabbitmq"
)

const (
	_defaultWaitTime = 2 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

type Client struct {
	conn *rmqrpc.Connection

	serverExchange string
	timeout        time.Duration

	error chan error
	stop  chan struct{}

	once sync.Once
}

func New(url, serverExchange, clientExchange string, topics []string, opts ...Option) (*Client, error) {
	const op = "rmq_client - New"

	cfg := rmqrpc.Params{
		URL:      url,
		WaitTime: _defaultWaitTime,
		Attempts: _defaultAttempts,
	}

	c := &Client{
		conn:           rmqrpc.NewConnection(clientExchange, topics, cfg),
		serverExchange: serverExchange,
		error:          make(chan error),
		stop:           make(chan struct{}),
		timeout:        _defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	err := c.conn.AttemptConnect(c.conn.ConnectWriter())
	if err != nil {
		return nil, fmt.Errorf("%s - c.conn.AttemptConnect: %w", op, err)
	}

	return c, nil
}

func (c *Client) publish(corrID, handler string, priority models.NotifyType, request any) error {
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

	err = c.conn.Channel.Publish(handler, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.conn.ConsumerExchange,
			Type:          handler,
			Body:          requestBody,
			Priority:      uint8(priority),
		})

	if err != nil {
		return fmt.Errorf("c.Channel.Publish: %w", err)
	}

	return nil
}

func (c *Client) reconnect() {
	close(c.stop)

	err := c.conn.AttemptConnect(c.conn.ConnectWriter())
	if err != nil {
		c.error <- err
		close(c.error)

		return
	}

	c.stop = make(chan struct{})
}

func (c *Client) RemoteCall(ctx context.Context, handler string, priority models.NotifyType, request any) error {
	const op = "rmq_client - RemoteCall"

	select {
	case <-c.stop:
		time.Sleep(c.timeout)
		select {
		case <-c.stop:
			return ErrConnectionClosed
		default:
		}
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	tracer := otel.Tracer("MessageClient")
	_, span := tracer.Start(ctx, "RemoteCall")
	defer span.End()

	corrID := uuid.New().String()
	span.SetAttributes(attribute.String("message.queue.id", corrID))

	err := c.publish(corrID, handler, priority, request)

	if err != nil {
		span.RecordError(err)

		unwrappedErr := errors.Unwrap(err)

		if amqpErr, ok := unwrappedErr.(*amqp.Error); ok {
			if amqpErr.Code == amqp.ChannelError || amqpErr.Code == amqp.ConnectionForced {
				go func() {
					c.once.Do(func() { c.reconnect() })
					time.Sleep(c.timeout)
					c.once = sync.Once{}
				}()
			}
		}

		return fmt.Errorf("%s - c.publish: %w", op, err)
	}
	return nil
}

func (c *Client) Notify() <-chan error {
	return c.error
}

func (c *Client) Shutdown() error {
	const op = "rmq_client - Shutdown"

	select {
	case <-c.error:
		return nil
	default:
	}

	close(c.stop)
	time.Sleep(c.timeout)

	err := c.conn.Connection.Close()
	if err != nil {
		return fmt.Errorf("%s - c.Connection.Close: %w", op, err)
	}

	return nil
}

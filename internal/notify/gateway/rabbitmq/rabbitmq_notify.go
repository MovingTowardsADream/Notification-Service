package rabbitmq

import (
	"Notification_Service/internal/entity"
	"context"
	"fmt"
)

type NotifyGatewayRMQ interface {
	RemoteCall(ctx context.Context, handler string, request interface{}) error
}

type NotifyGateway struct {
	rmq NotifyGatewayRMQ
}

func New(rmq NotifyGatewayRMQ) *NotifyGateway {
	return &NotifyGateway{rmq}
}

func (gw *NotifyGateway) CreateNotifyMessageOnRabbitMQ(ctx context.Context, notify entity.Notify) error {
	err := wrapper(ctx, func() error {
		return gw.rmq.RemoteCall(ctx, "createNewNotify", notify)
	})

	if err != nil {
		return fmt.Errorf("NotifyGateway - CreateNotifyMessageOnRabbitMQ - gw.rmq.RemoteCall: %w", err)
	}

	return nil
}

func wrapper(ctx context.Context, f func() error) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

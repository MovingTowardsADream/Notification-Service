package gotests

import (
	"context"

	amqprpc "Notification_Service/internal/infrastructure/workers/amqp_rpc"
	"Notification_Service/tests/gotests/mocks"
)

func SetupMocks(ctx context.Context, name string) (
	context.Context,
	Repository,
	func(),
) {
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, TestSessionIDHeader, name)

	rmqRouter := amqprpc.NewRouter(new(mocks.MockSenderMail), new(mocks.MockSenderPhone))

	repo, err := NewRepository(ctx, rmqRouter)
	if err != nil {
		panic(err)
	}

	err = repo.Start(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, repo, func() {
		_ = repo.Stop(context.Background())
		cancel()
	}
}

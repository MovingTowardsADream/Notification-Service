package gotests

import (
	"context"

	"Notification_Service/pkg/logger"
)

func SetupMocks(ctx context.Context, name string) (
	context.Context,
	*MockServer,
	Repository,
	func(),
) {
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, TestSessionIDHeader, name)

	repo, err := NewRepository()
	if err != nil {
		panic(err)
	}

	err = repo.Start(ctx)
	if err != nil {
		panic(err)
	}

	cfg := repo.Config()

	log, err := logger.Setup(cfg.Log.Level, cfg.Log.Path)
	if err != nil {
		panic(err)
	}

	mockServer := NewMockServer(log, , WithPort(8081))

	return ctx, mockServer, repo, func() {
		_ = mockServer.Close()
		_ = repo.Stop(context.Background())
		cancel()
	}
}

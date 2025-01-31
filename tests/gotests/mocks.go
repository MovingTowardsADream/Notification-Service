package gotests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"Notification_Service/internal/application/usecase"
	gwmessaging "Notification_Service/internal/infrastructure/gateway/messaging"
	"Notification_Service/internal/infrastructure/repository/postgres"
	amqprpc "Notification_Service/internal/infrastructure/workers/amqp_rpc"
	logMocks "Notification_Service/pkg/logger/mocks"
	"Notification_Service/tests/gotests/mocks"
)

func SetupMocks(ctx context.Context, name string, t *testing.T) (
	context.Context,
	Repository,
	*MockServer,
	func(),
) {
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, TestSessionIDHeader, name)

	ctrl := gomock.NewController(t)

	mailMocks := mocks.NewMockSenderMail(ctrl)
	mailMocks.EXPECT().SendMailLetter(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	phoneMocks := mocks.NewMockSenderPhone(ctrl)
	phoneMocks.EXPECT().SendPhoneSMS(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	logger := logMocks.NewMockLogger(ctrl)
	logger.EXPECT().Error(gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Warn(gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any()).AnyTimes()

	rmqRouter := amqprpc.NewRouter(mailMocks, phoneMocks)

	repo, err := NewRepository(ctx, rmqRouter)
	if err != nil {
		panic(err)
	}

	usersRepo := postgres.NewUsersRepo(repo.Storage())

	gateway := gwmessaging.NewNotifyGateway(repo.MesClient())

	notifySender := usecase.NewNotifySender(logger, usersRepo, gateway)
	editInfo := usecase.NewEditInfo(logger, usersRepo)

	mockServer := NewMockServer(notifySender, editInfo)

	err = repo.Start(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, repo, mockServer, func() {
		ctrl.Finish()
		_ = repo.Stop(context.Background())
		_ = mockServer.Close()
		cancel()
	}
}

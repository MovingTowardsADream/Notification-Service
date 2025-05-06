package gotests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"Notification_Service/internal/application/usecase"
	gwmessaging "Notification_Service/internal/infrastructure/gateway/messaging"
	"Notification_Service/internal/infrastructure/repository/postgres/notify"
	"Notification_Service/internal/infrastructure/repository/postgres/users"
	amqprpc "Notification_Service/internal/infrastructure/workers/amqp_rpc"
	"Notification_Service/internal/infrastructure/workers/outbox"
	outboxhandler "Notification_Service/internal/infrastructure/workers/outbox/handlers"
	"Notification_Service/pkg/hasher"
	"Notification_Service/tests/gotests/mocks"
)

func SetupMocks(ctx context.Context, name string, t *testing.T) (
	context.Context,
	Repository,
	*MockServer,
	StackCancelFunc[CancelFunc],
) {
	stackCancel := make(StackCancelFunc[CancelFunc], 0)

	defer func() {
		if err := recover(); err != nil {
			_ = stackCancel.Clear()
			panic(err)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	stackCancel.Push(func() { cancel() })
	ctx = context.WithValue(ctx, TestSessionIDHeader, name)

	ctrl := gomock.NewController(t)
	stackCancel.Push(func() { ctrl.Finish() })

	mailMocks := mocks.NewMockSenderMail(ctrl)
	mailMocks.EXPECT().SendMailLetter(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	phoneMocks := mocks.NewMockSenderPhone(ctrl)
	phoneMocks.EXPECT().SendPhoneSMS(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	rmqRouter := amqprpc.NewRouter(mailMocks, phoneMocks)

	repo, err := NewRepository(ctx, rmqRouter)
	if err != nil {
		panic(err)
	}

	stackCancel.Push(func() { _ = repo.Stop(context.Background()) })

	err = repo.Start(ctx)
	if err != nil {
		panic(err)
	}

	usersRepo := users.NewUsersRepo(repo.Storage())
	notifyRepo := notify.NewNotifyRepo(repo.Storage())

	gateway := gwmessaging.NewNotifyGateway(repo.MesClient())

	outboxWorker := outbox.NewWorker(
		map[string]outbox.WorkerRun{
			"mail":  outboxhandler.NewMailWorker(repo.Logger(), notifyRepo, gateway),
			"phone": outboxhandler.NewPhoneWorker(repo.Logger(), notifyRepo, gateway),
		},
	)

	err = outboxWorker.WorkerRun()
	if err != nil {
		panic(err)
	}

	stackCancel.Push(func() { _ = outboxWorker.Shutdown() })

	notifySender := usecase.NewNotifySender(repo.Logger(), usersRepo, notifyRepo)

	hash := hasher.NewSHA1Hasher(repo.Config().Security.PasswordSalt)

	editInfo := usecase.NewUserInfo(repo.Logger(), usersRepo, hash)

	mockServer := NewMockServer(notifySender, editInfo)

	err = mockServer.ListenAndServe(ctx)
	if err != nil {
		panic(err)
	}

	stackCancel.Push(func() { _ = mockServer.Close() })

	return ctx, repo, mockServer, stackCancel
}

package gotests

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/infrastructure/config"
	"Notification_Service/pkg/utils"
)

type Clients struct {
	Notify notifyv1.NotifyClient
	Users  notifyv1.UsersClient
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	host, port := "", cfg.GRPC.Port

	nc, err := grpc.DialContext(ctx, utils.FormatAddress(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	uc, err := grpc.DialContext(ctx, utils.FormatAddress(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Clients{
		Notify: notifyv1.NewNotifyClient(nc),
		Users:  notifyv1.NewUsersClient(uc),
	}, nil
}

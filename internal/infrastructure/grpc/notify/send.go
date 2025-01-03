package notify

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/application/usecase"
	grpc_err "Notification_Service/internal/infrastructure/grpc/errors"
	"Notification_Service/internal/interfaces/dto"
)

type NotifySend interface {
	SendNotifyForUser(ctx context.Context, notifyRequest *dto.ReqNotification) error
}

type sendNotifyRoutes struct {
	notifyv1.UnimplementedNotifyServer
	notifySend NotifySend
}

func Notify(gRPC *grpc.Server, notifySend NotifySend) {
	notifyv1.RegisterNotifyServer(gRPC, &sendNotifyRoutes{notifySend: notifySend})
}

func (s *sendNotifyRoutes) SendMessage(ctx context.Context, req *notifyv1.SendMessageReq) (*notifyv1.SendMessageResp, error) {
	requestNotification := &dto.ReqNotification{
		UserID:     req.UserID,
		NotifyType: req.NotifyType.String(),
		Channels: dto.Channels{
			Mail: &dto.MailChannel{
				Subject: req.Channels.Mail.Subject,
				Body:    req.Channels.Mail.Body,
			},
			Phone: &dto.PhoneChannel{
				Body: req.Channels.Phone.Body,
			},
		},
	}

	err := s.notifySend.SendNotifyForUser(ctx, requestNotification)

	if err != nil {
		if errors.Is(err, usecase.ErrTimeout) {
			return nil, grpc_err.ErrDeadlineExceeded
		} else if errors.Is(err, usecase.ErrNotFound) {
			return nil, grpc_err.ErrNotFound
		}

		return nil, grpc_err.ErrInternalServer
	}

	return &notifyv1.SendMessageResp{
		Respond: "success",
	}, nil
}

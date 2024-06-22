package grpc_send_notify

import (
	"Notification_Service/internal/entity"
	grpc_error "Notification_Service/internal/notify/grpc/error"
	"Notification_Service/internal/notify/usecase/usecase_errors"
	custom_validator "Notification_Service/pkg/validator"
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"errors"
	"google.golang.org/grpc"
)

type sendNotifyRoutes struct {
	notifyv1.UnimplementedSendNotifyServer
	notifySend NotifySend
	validator  *custom_validator.CustomValidator
}

func SendNotify(gRPC *grpc.Server, notifySend NotifySend, v *custom_validator.CustomValidator) {
	notifyv1.RegisterSendNotifyServer(gRPC, &sendNotifyRoutes{notifySend: notifySend, validator: v})
}

func (s *sendNotifyRoutes) SendMessage(ctx context.Context, req *notifyv1.SendMessageRequest) (*notifyv1.SendMessageResponse, error) {

	if err := s.validator.Validate(req); err != nil {
		return nil, grpc_error.ErrInvalidArgument(err)
	}

	requestNotification := &entity.RequestNotification{
		UserId:     req.UserId,
		NotifyType: req.NotifyType.String(),
		Channels: entity.Channels{
			Mail: entity.MailChannel{
				Subject: req.Channels.Mail.Subject,
				Body:    req.Channels.Mail.Body,
			},
			Phone: entity.PhoneChannel{
				Body: req.Channels.Phone.Body,
			},
		},
	}

	err := s.notifySend.SendNotifyForUser(ctx, requestNotification)

	if err != nil {
		if errors.Is(err, usecase_errors.ErrTimeout) {
			return nil, grpc_error.ErrDeadlineExceeded
		} else if errors.Is(err, usecase_errors.ErrNotFound) {
			return nil, grpc_error.ErrNotFound
		}

		return nil, grpc_error.ErrInternalServer
	}

	return &notifyv1.SendMessageResponse{
		Respond: "success",
	}, nil
}

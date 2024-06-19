package grpc_send_notify

import (
	"Notification_Service/internal/entity"
	grpc_error "Notification_Service/internal/notify/grpc/error"
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"errors"
	"google.golang.org/grpc"
)

type sendNotifyRoutes struct {
	notifyv1.UnimplementedSendNotifyServer
	notifySend NotifySend
}

func SendNotify(gRPC *grpc.Server, notifySend NotifySend) {
	notifyv1.RegisterSendNotifyServer(gRPC, &sendNotifyRoutes{notifySend: notifySend})
}

func (s *sendNotifyRoutes) SendMessage(ctx context.Context, req *notifyv1.SendMessageRequest) (*notifyv1.SendMessageResponse, error) {

	// TODO: Validate request

	err := s.notifySend.SendNotifyForUser(ctx, req)

	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			return nil, grpc_error.ErrDeadlineExceeded
		}

		// TODO logging error

		return nil, grpc_error.ErrInternalServer

	}

	return &notifyv1.SendMessageResponse{
		Respond: "success",
	}, nil
}

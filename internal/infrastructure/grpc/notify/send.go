package notify

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/application/usecase"
	grpcerr "Notification_Service/internal/infrastructure/grpc/errors"
	"Notification_Service/internal/interfaces/convert"
	"Notification_Service/internal/interfaces/dto"
)

type SendersNotify interface {
	SendToUser(ctx context.Context, notifyRequest *dto.ReqNotification) error
}

type sendNotifyRoutes struct {
	notifyv1.UnimplementedNotifyServer
	notifySend SendersNotify
}

func Notify(gRPC *grpc.Server, notifySend SendersNotify) {
	notifyv1.RegisterNotifyServer(gRPC, &sendNotifyRoutes{notifySend: notifySend})
}

func (s *sendNotifyRoutes) SendMessage(ctx context.Context, req *notifyv1.SendMessageReq) (*notifyv1.SendMessageResp, error) {
	tracer := otel.Tracer("userRoutes")
	ctx, span := tracer.Start(ctx, "SendMessage")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", req.GetUserID()))

	if err := req.ValidateAll(); err != nil {
		span.RecordError(err)
		return nil, grpcerr.ErrInvalidArgument
	}

	dataNotify, err := convert.SendMessageReqToReqNotification(req)

	if err != nil {
		span.RecordError(err)
		return nil, grpcerr.ErrInvalidArgument
	}

	err = s.notifySend.SendToUser(ctx, dataNotify)

	if err != nil {
		span.RecordError(err)

		if errors.Is(err, usecase.ErrTimeout) {
			return nil, grpcerr.ErrDeadlineExceeded
		} else if errors.Is(err, usecase.ErrNotFound) {
			return nil, grpcerr.ErrNotFound
		}

		return nil, grpcerr.ErrInternalServer
	}

	return &notifyv1.SendMessageResp{
		Respond: "success",
	}, nil
}

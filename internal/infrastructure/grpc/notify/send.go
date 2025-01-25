package notify

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
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
	const (
		tracerName = "userRoutes"
		spanName   = "SendMessage"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
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

		return nil, grpcerr.MappingErrors(err)
	}

	return &notifyv1.SendMessageResp{
		Respond: "success",
	}, nil
}

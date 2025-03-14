package notify

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	notifyv1 "Notification_Service/api/gen/go/notify"
	grpcerr "Notification_Service/internal/infrastructure/grpc/errors"
	"Notification_Service/internal/interfaces/convert"
	"Notification_Service/internal/interfaces/dto"
)

type SendersNotify interface {
	SendToUser(ctx context.Context, notifyRequest *dto.ReqNotification) error
}

func (s *sendNotifyRoutes) SendMessage(ctx context.Context, req *notifyv1.SendMessageReq) (*notifyv1.SendMessageResp, error) {
	const (
		tracerName = "userRoutes"
		spanName   = "SendMessage"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	if err := req.ValidateAll(); err != nil {
		span.RecordError(err)
		return nil, grpcerr.ErrInvalidArgument
	}

	dataNotify, err := convert.SendMessageReqToReqNotification(ctx, req)
	if err != nil {
		span.RecordError(err)
		return nil, grpcerr.ErrInvalidArgument
	}

	span.SetAttributes(attribute.String("user.id", dataNotify.UserID))
	span.SetAttributes(attribute.String("request.id", dataNotify.RequestID))

	err = s.notifySend.SendToUser(ctx, dataNotify)

	if err != nil {
		span.RecordError(err)

		return nil, grpcerr.MappingErrors(err)
	}

	return &notifyv1.SendMessageResp{
		Respond: "success",
	}, nil
}

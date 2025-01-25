package users

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

type EditInfo interface {
	EditUserPreferences(ctx context.Context, preferences *dto.UserPreferences) error
}

type userRoutes struct {
	notifyv1.UnimplementedUsersServer
	editInfo EditInfo
}

func Users(gRPC *grpc.Server, editInfo EditInfo) {
	notifyv1.RegisterUsersServer(gRPC, &userRoutes{editInfo: editInfo})
}

func (s *userRoutes) EditPreferences(ctx context.Context, req *notifyv1.EditPreferencesReq) (*notifyv1.EditPreferencesResp, error) {
	const (
		tracerName = "userRoutes"
		spanName   = "EditPreferences"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.String("user.id", req.GetUserID()))

	if err := req.ValidateAll(); err != nil {
		span.RecordError(err)
		return nil, grpcerr.ErrInvalidArgument
	}

	preferences := convert.EditPreferencesReqToUserPreferences(req)

	err := s.editInfo.EditUserPreferences(ctx, preferences)

	if err != nil {
		span.RecordError(err)

		return nil, grpcerr.MappingErrors(err)
	}

	return &notifyv1.EditPreferencesResp{
		Respond: "success",
	}, nil
}

package users

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/timestamppb"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/domain/models"
	grpcerr "Notification_Service/internal/infrastructure/grpc/errors"
	"Notification_Service/internal/interfaces/convert"
	"Notification_Service/internal/interfaces/dto"
)

type UserInfo interface {
	AddUser(ctx context.Context, user *dto.User) (*models.User, error)
	EditUserPreferences(ctx context.Context, preferences *dto.UserPreferences) error
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

	err := s.userInfo.EditUserPreferences(ctx, preferences)

	if err != nil {
		span.RecordError(err)

		return nil, grpcerr.MappingErrors(err)
	}

	return &notifyv1.EditPreferencesResp{
		Respond: "success",
	}, nil
}

func (s *userRoutes) AddUser(ctx context.Context, req *notifyv1.AddUserReq) (*notifyv1.AddUserResp, error) {
	const (
		tracerName = "userRoutes"
		spanName   = "AddUser"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	if err := req.ValidateAll(); err != nil {
		span.RecordError(err)
		return nil, grpcerr.ErrInvalidArgument
	}

	userReq := convert.AddUserReqToUser(req)

	user, err := s.userInfo.AddUser(ctx, userReq)

	if err != nil {
		span.RecordError(err)

		return nil, grpcerr.MappingErrors(err)
	}

	return &notifyv1.AddUserResp{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

package users

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/application/usecase"
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
	if err := req.ValidateAll(); err != nil {
		return nil, grpcerr.ErrInvalidArgument
	}

	preferences := convert.EditPreferencesReqToUserPreferences(req)

	err := s.editInfo.EditUserPreferences(ctx, preferences)

	if err != nil {
		if errors.Is(err, usecase.ErrTimeout) {
			return nil, grpcerr.ErrDeadlineExceeded
		} else if errors.Is(err, usecase.ErrNotFound) {
			return nil, grpcerr.ErrNotFound
		}

		return nil, grpcerr.ErrInternalServer
	}

	return &notifyv1.EditPreferencesResp{
		Respond: "success",
	}, nil
}

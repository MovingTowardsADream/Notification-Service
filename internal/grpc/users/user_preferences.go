package grpc_users

import (
	"Notification_Service/internal/entity"
	grpc_error "Notification_Service/internal/grpc/error"
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"errors"
	"google.golang.org/grpc"
)

type EditInfo interface {
	EditPreferences(ctx context.Context, preferences *notifyv1.UserPreferencesRequest) error
}

type userRoutes struct {
	notifyv1.UnimplementedUsersServer
	editInfo EditInfo
}

func Users(gRPC *grpc.Server, editInfo EditInfo) {
	notifyv1.RegisterUsersServer(gRPC, &userRoutes{editInfo: editInfo})
}

func (s *userRoutes) UserPreferences(ctx context.Context, req *notifyv1.UserPreferencesRequest) (*notifyv1.UserPreferencesResponse, error) {

	// TODO: Validate request

	err := s.editInfo.EditPreferences(ctx, req)

	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			return nil, grpc_error.ErrDeadlineExceeded
		}

		// TODO logging error

		return nil, grpc_error.ErrInternalServer

	}

	return &notifyv1.UserPreferencesResponse{
		Respond: "success",
	}, nil
}

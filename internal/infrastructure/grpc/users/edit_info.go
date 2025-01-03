package users

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/application/usecase"
	grpc_err "Notification_Service/internal/infrastructure/grpc/errors"
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
	preferences := &dto.UserPreferences{
		UserID:      req.UserID,
		Preferences: dto.Preferences{},
	}

	if req.Preferences.Mail != nil {
		preferences.Preferences.Mail = &dto.MailPreference{
			Approval: req.Preferences.Mail.Approval,
		}
	}

	if req.Preferences.Phone != nil {
		preferences.Preferences.Phone = &dto.PhonePreference{
			Approval: req.Preferences.Phone.Approval,
		}
	}

	err := s.editInfo.EditUserPreferences(ctx, preferences)

	if err != nil {
		if errors.Is(err, usecase.ErrTimeout) {
			return nil, grpc_err.ErrDeadlineExceeded
		} else if errors.Is(err, usecase.ErrNotFound) {
			return nil, grpc_err.ErrNotFound
		}

		return nil, grpc_err.ErrInternalServer
	}

	return &notifyv1.EditPreferencesResp{
		Respond: "success",
	}, nil
}

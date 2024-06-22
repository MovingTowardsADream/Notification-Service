package grpc_users

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notify/grpc/error"
	"Notification_Service/internal/notify/usecase/usecase_errors"
	custom_validator "Notification_Service/pkg/validator"
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"errors"
	"google.golang.org/grpc"
)

type userRoutes struct {
	notifyv1.UnimplementedUsersServer
	editInfo  EditInfo
	validator *custom_validator.CustomValidator
}

func Users(gRPC *grpc.Server, editInfo EditInfo, v *custom_validator.CustomValidator) {
	notifyv1.RegisterUsersServer(gRPC, &userRoutes{editInfo: editInfo, validator: v})
}

func (s *userRoutes) UserPreferences(ctx context.Context, req *notifyv1.UserPreferencesRequest) (*notifyv1.UserPreferencesResponse, error) {

	if err := s.validator.Validate(req); err != nil {
		return nil, grpc_error.ErrInvalidArgument(err)
	}

	preferences := &entity.UserPreferences{
		UserId:      req.UserId,
		Preferences: entity.Preferences{},
	}

	if req.Preferences.Mail != nil {
		preferences.Preferences.Mail = &entity.MailPreference{
			Approval: req.Preferences.Mail.Approval,
		}
	}

	if req.Preferences.Phone != nil {
		preferences.Preferences.Phone = &entity.PhonePreference{
			Approval: req.Preferences.Phone.Approval,
		}
	}

	err := s.editInfo.EditUserPreferences(ctx, preferences)

	if err != nil {
		if errors.Is(err, usecase_errors.ErrTimeout) {
			return nil, grpc_error.ErrDeadlineExceeded
		} else if errors.Is(err, usecase_errors.ErrNotFound) {
			return nil, grpc_error.ErrNotFound
		}

		return nil, grpc_error.ErrInternalServer

	}

	return &notifyv1.UserPreferencesResponse{
		Respond: "success",
	}, nil
}

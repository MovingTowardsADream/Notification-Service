package usecase

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/domain/models"
	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/hasher"
	"Notification_Service/pkg/logger"
)

type Users interface {
	Create(ctx context.Context, user *dto.User) (*models.User, error)
	EditPreferences(ctx context.Context, preferences *dto.UserPreferences) error
}

type UserInfo struct {
	l         logger.Logger
	usersData Users
	hash      hasher.PasswordHash
}

func NewUserInfo(l logger.Logger, usersDataPref Users, hash hasher.PasswordHash) *UserInfo {
	return &UserInfo{
		l:         l,
		usersData: usersDataPref,
		hash:      hash,
	}
}

func (e *UserInfo) EditUserPreferences(ctx context.Context, preferences *dto.UserPreferences) error {
	const op = "UserInfo - EditUserPreferences"

	const (
		tracerName = "UserInfo"
		spanName   = "EditUserPreferences"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.String("user.id", preferences.UserID))

	err := e.usersData.EditPreferences(ctx, preferences)

	if err != nil {
		span.RecordError(err)

		if ok, err := mappingErrors(err); ok {
			return err
		}

		e.l.Error(
			op,
			e.l.Err(err),
			logger.AnyAttr("trace-id", ctx.Value("trace-id").(string)),
		)

		return err
	}

	return nil
}

func (e *UserInfo) AddUser(ctx context.Context, userData *dto.User) (*models.User, error) {
	const op = "UserInfo - AddUser"

	const (
		tracerName = "EditInfo"
		spanName   = "EditUserPreferences"
	)

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	userData.Password = e.hash.Hash(userData.Password)

	user, err := e.usersData.Create(ctx, userData)

	if err != nil {
		span.RecordError(err)

		if ok, err := mappingErrors(err); ok {
			return nil, err
		}

		e.l.Error(
			op,
			e.l.Err(err),
			logger.AnyAttr("trace-id", ctx.Value("trace-id").(string)),
		)

		return nil, err
	}

	return user, nil
}

package usecase

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	repoerr "Notification_Service/internal/infrastructure/repository/errors"
	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

const (
	_defaultTimeout = 5 * time.Second
)

type UserPreferences interface {
	EditPreferences(ctx context.Context, preferences *dto.UserPreferences) error
}

type EditInfo struct {
	l         *logger.Logger
	usersData UserPreferences
}

func NewEditInfo(l *logger.Logger, usersDataPref UserPreferences) *EditInfo {
	return &EditInfo{
		l:         l,
		usersData: usersDataPref,
	}
}

func (e *EditInfo) EditUserPreferences(ctx context.Context, preferences *dto.UserPreferences) error {
	tracer := otel.Tracer("EditInfo")
	ctx, span := tracer.Start(ctx, "EditUserPreferences")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", preferences.UserID))

	err := e.usersData.EditPreferences(ctx, preferences)

	if err != nil {
		span.RecordError(err)

		if errors.Is(err, repoerr.ErrNotFound) {
			return ErrNotFound
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}

		e.l.Error(
			"EditUserPreferences - e.usersDataPref.EditPreferences",
			e.l.Err(err),
			logger.NewStrArgs("trace-id", ctx.Value("trace-id").(string)),
		)

		return err
	}

	return nil
}

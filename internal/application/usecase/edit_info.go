package usecase

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

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
	const op = "EditInfo - EditUserPreferences"

	const (
		tracerName = "EditInfo"
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
			logger.NewStrArgs("trace-id", ctx.Value("trace-id").(string)),
		)

		return err
	}

	return nil
}

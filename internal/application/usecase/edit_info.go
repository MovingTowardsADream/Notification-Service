package usecase

import (
	"context"
	"errors"
	"time"

	repoErr "Notification_Service/internal/infrastructure/repository/errors"
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
	err := e.usersData.EditPreferences(ctx, preferences)

	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return ErrNotFound
		} else if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}

		e.l.Error("EditUserPreferences - e.usersDataPref.EditPreferences: ", e.l.Err(err))

		return err
	}

	return nil
}

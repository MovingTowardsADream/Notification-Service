package usecase

import (
	"context"
	"errors"
	"time"

	"Notification_Service/internal/interfaces/dto"
	"Notification_Service/pkg/logger"
)

const (
	_defaultTimeout = 5 * time.Second
)

type UsersDataPreferences interface {
	EditUserPreferences(ctx context.Context, preferences *dto.UserPreferences) error
}

type EditInfo struct {
	l             *logger.Logger
	usersDataPref UsersDataPreferences
}

func NewEditInfo(l *logger.Logger, usersDataPref UsersDataPreferences) *EditInfo {
	return &EditInfo{
		l:             l,
		usersDataPref: usersDataPref,
	}
}

func (e *EditInfo) EditUserPreferences(ctx context.Context, preferences *dto.UserPreferences) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	err := e.usersDataPref.EditUserPreferences(ctxTimeout, preferences)

	if err != nil {
		if errors.Is(err, ErrTimeout) {
			return ErrTimeout
		} else if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		e.l.Error("EditUserPreferences - e.usersDataPref.EditUserPreferences: ", e.l.Err(err))

		return err
	}

	return nil
}

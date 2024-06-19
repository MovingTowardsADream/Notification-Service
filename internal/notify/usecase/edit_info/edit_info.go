package edit_info

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notify/repository/repository_erros"
	"Notification_Service/internal/notify/usecase/usecase_errors"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

const (
	_defaultTimeout = 5 * time.Second
)

type EditInfo struct {
	l             *slog.Logger
	usersDataPref UsersDataPreferences
}

func New(l *slog.Logger, usersDataPref UsersDataPreferences) *EditInfo {
	return &EditInfo{
		l:             l,
		usersDataPref: usersDataPref,
	}
}

func (e *EditInfo) EditUserPreferences(ctx context.Context, preferences *entity.UserPreferences) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	err := e.usersDataPref.EditUserPreferences(ctxTimeout, preferences)

	if err != nil {
		if errors.Is(err, usecase_errors.ErrTimeout) {
			return usecase_errors.ErrTimeout
		} else if errors.Is(err, repository_erros.ErrNotFound) {
			return usecase_errors.ErrNotFound
		}

		// TODO logging error

		return fmt.Errorf("UseCase - EditInfo - EditUserPreferences - e.usersDataPref.EditUserPreferences: %w", err)
	}

	return nil
}

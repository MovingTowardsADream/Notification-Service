package edit_info

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notify/repository/repository_erros"
	"Notification_Service/internal/notify/usecase/usecase_errors"
	"Notification_Service/pkg/logger"
	"context"
	"errors"
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

		e.l.Error("EditUserPreferences - e.usersDataPref.EditUserPreferences: ", logger.Err(err))

		return err
	}

	return nil
}

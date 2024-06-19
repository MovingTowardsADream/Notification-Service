package edit_info

import (
	"Notification_Service/internal/entity"
	"context"
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

	// TODO logging

	return e.usersDataPref.EditUserPreferences(ctxTimeout, preferences)
}

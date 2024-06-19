package edit_info

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"log/slog"
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

func (e *EditInfo) EditUserPreferences(ctx context.Context, preferences *notifyv1.UserPreferencesRequest) error {

	_ = e.usersDataPref.EditUserPreferences(ctx, preferences)

	// TODO logging

	// TODO other logic

	return nil
}

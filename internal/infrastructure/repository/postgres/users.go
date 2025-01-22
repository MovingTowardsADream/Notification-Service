package postgres

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/interfaces/dto"
)

const (
	usersTable  = "users"
	notifyTable = "notifications"
)

type UsersRepo struct {
	storage *Postgres
}

func NewUsersRepo(storage *Postgres) *UsersRepo {
	return &UsersRepo{storage: storage}
}

func (ur *UsersRepo) GetUserCommunication(
	ctx context.Context,
	communication *dto.IdentificationUserCommunication,
) (*dto.UserCommunication, error) {
	tracer := otel.Tracer("UsersRepo")
	ctx, span := tracer.Start(ctx, "GetUserCommunication")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", communication.ID))

	sql, args, _ := ur.storage.Builder.
		Select("users.id", "users.email", "users.phone", "notifications.email_notify", "notifications.phone_notify").
		From(usersTable).
		InnerJoin("notifications on users.id = notifications.user_id").
		Where("users.id = ?", communication.ID).
		ToSql()

	var userCommunication dto.UserCommunication

	err := ur.storage.Pool.QueryRow(ctx, sql, args...).Scan(
		&userCommunication.ID,
		&userCommunication.Email,
		&userCommunication.Phone,
		&userCommunication.MailPref,
		&userCommunication.PhonePref,
	)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("NotifyRepo.GetUserCommunication - r.Pool.QueryRow: %w", mappingErrors(err))
	}

	return &userCommunication, nil
}

func (ur *UsersRepo) EditPreferences(ctx context.Context, preferences *dto.UserPreferences) error {
	tracer := otel.Tracer("UsersRepo")
	ctx, span := tracer.Start(ctx, "EditPreferences")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", preferences.UserID))

	tx, err := ur.storage.Pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("UsersRepo.EditPreferences - r.Pool.Begin: %w", mappingErrors(err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := ur.storage.Builder.
		Select("email_notify", "phone_notify").
		From(notifyTable).
		Where("user_id = ?", preferences.UserID).
		ToSql()

	var emailNotify, phoneNotify bool
	err = tx.QueryRow(ctx, sql, args...).Scan(&emailNotify, &phoneNotify)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("UsersRepo.EditPreferences - tx.QueryRow: %w", mappingErrors(err))
	}

	if preferences.Preferences.Mail == nil {
		preferences.Preferences.Mail = &dto.MailPreference{
			Approval: emailNotify,
		}
	}
	if preferences.Preferences.Phone == nil {
		preferences.Preferences.Phone = &dto.PhonePreference{
			Approval: phoneNotify,
		}
	}

	sql, args, _ = ur.storage.Builder.
		Update(notifyTable).
		Set("email_notify", preferences.Preferences.Mail.Approval).
		Set("phone_notify", preferences.Preferences.Phone.Approval).
		Where("user_id = ?", preferences.UserID).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("UsersRepo.EditPreferences - tx.Exec: %w", mappingErrors(err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("UsersRepo.EditPreferences - tx.Commit: %w", mappingErrors(err))
	}

	return nil
}

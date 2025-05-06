package users

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/domain/models"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/interfaces/dto"
)

const (
	usersTable  = "users"
	notifyTable = "notifications"
)

const tracerName = "userRepo"

type RepoUsers struct {
	storage *postgres.Postgres
}

func NewUsersRepo(storage *postgres.Postgres) *RepoUsers {
	return &RepoUsers{storage: storage}
}

func (ur *RepoUsers) GetUserCommunication(
	ctx context.Context,
	communication *dto.IdentificationUserCommunication,
) (*dto.UserCommunication, error) {
	const op = "NotifyRepo.GetUserCommunication"
	const spanName = "GetUserCommunication"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.String("user.id", communication.ID))

	sql, args, _ := ur.storage.Builder.
		Select("users.id", "notifications.email_notify", "notifications.phone_notify").
		From(usersTable).
		InnerJoin("notifications on users.id = notifications.user_id").
		Where("users.id = ?", communication.ID).
		ToSql()

	var userCommunication dto.UserCommunication

	err := ur.storage.Pool.QueryRow(ctx, sql, args...).Scan(
		&userCommunication.ID,
		&userCommunication.MailPref,
		&userCommunication.PhonePref,
	)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("%s - r.Pool.QueryRow: %w", op, postgres.MappingErrors(err))
	}

	return &userCommunication, nil
}

func (ur *RepoUsers) EditPreferences(ctx context.Context, preferences *dto.UserPreferences) error {
	const op = "UsersRepo.EditPreferences"
	const spanName = "EditPreferences"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.String("user.id", preferences.UserID))

	tx, err := ur.storage.Pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - r.Pool.Begin: %w", op, postgres.MappingErrors(err))
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
		return fmt.Errorf("%s - tx.QueryRow: %w", op, postgres.MappingErrors(err))
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
		return fmt.Errorf("%s - tx.Exec: %w", op, postgres.MappingErrors(err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - tx.Commit: %w", op, postgres.MappingErrors(err))
	}

	return nil
}

func (ur *RepoUsers) Create(ctx context.Context, userData *dto.User) (*models.User, error) {
	const op = "UsersRepo.Add"
	const spanName = "Add"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	tx, err := ur.storage.Pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("%s - r.Pool.Begin: %w", op, postgres.MappingErrors(err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	user := &models.User{}

	sql, args, _ := ur.storage.Builder.
		Insert(usersTable).
		Columns("username", "email", "phone", "password_hash").
		Values(userData.Username, userData.Email, userData.Phone, userData.Password).
		Suffix("RETURNING id, username, email, phone, time").
		ToSql()

	err = ur.storage.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.CreatedAt)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("UserRepo.Create - ur.storage.Pool.QueryRow: %w", postgres.MappingErrors(err))
	}

	var emailNotify, phoneNotify bool

	if userData.Preferences.Mail != nil {
		emailNotify = userData.Preferences.Mail.Approval
	}

	if userData.Preferences.Phone != nil {
		phoneNotify = userData.Preferences.Phone.Approval
	}

	sql, args, _ = ur.storage.Builder.
		Insert(notifyTable).
		Columns("email_notify", "phone_notify", "user_id").
		Values(emailNotify, phoneNotify, user.ID).
		ToSql()

	_, err = ur.storage.Pool.Exec(ctx, sql, args...)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("UserRepo.Create - ur.storage.Pool.QueryRow: %w", postgres.MappingErrors(err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("%s - tx.Commit: %w", op, postgres.MappingErrors(err))
	}

	return user, nil
}

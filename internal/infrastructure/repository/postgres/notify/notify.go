package notify

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/interfaces/dto"
)

const (
	mailHistoryTable  = "history_email_notify"
	phoneHistoryTable = "history_phone_notify"
)

const tracerName = "notifyRepo"

type NotifyRepo struct {
	storage *postgres.Postgres
}

func NewNotifyRepo(storage *postgres.Postgres) *NotifyRepo {
	return &NotifyRepo{storage: storage}
}

func (nr *NotifyRepo) ProcessedNotify(
	ctx context.Context,
	processed *dto.ProcessedNotify,
) error {
	const op = "NotifyRepo.ProcessedNotify"
	const spanName = "ProcessedNotify"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.String("user.id", processed.UserID))
	span.SetAttributes(attribute.String("request.id", processed.RequestID))

	tx, err := nr.storage.Pool.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - r.Pool.Begin: %w", op, postgres.MappingErrors(err))
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if processed.MailDate != nil {
		sql, args, _ := nr.storage.Builder.
			Insert(mailHistoryTable).
			Columns("request_id", "notify_type", "subject", "body", "status", "user_id").
			Values(processed.RequestID, processed.MailDate.NotifyType, processed.MailDate.Subject, processed.MailDate.Body, "processed", processed.UserID).
			ToSql()

		_, err = nr.storage.Pool.Exec(ctx, sql, args...)

		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("%s - nr.storage.Pool.Exec: %w", op, postgres.MappingErrors(err))
		}
	}

	if processed.PhoneDate != nil {
		sql, args, _ := nr.storage.Builder.
			Insert(phoneHistoryTable).
			Columns("request_id", "notify_type", "body", "status", "user_id").
			Values(processed.RequestID, processed.PhoneDate.NotifyType, processed.PhoneDate.Body, "processed", processed.UserID).
			ToSql()

		_, err = nr.storage.Pool.Exec(ctx, sql, args...)

		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("%s - nr.storage.Pool.Exec: %w", op, postgres.MappingErrors(err))
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - tx.Commit: %w", op, postgres.MappingErrors(err))
	}

	return nil
}

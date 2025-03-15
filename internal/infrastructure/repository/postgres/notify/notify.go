package notify

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"Notification_Service/internal/domain/models"
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

func (nr *NotifyRepo) GetBatchMailNotify(ctx context.Context, batch *dto.BatchNotify) ([]*dto.MailIdempotencyData, error) {
	const op = "NotifyRepo.GetBatchMailNotify"
	const spanName = "GetBatchMailNotify"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	sql, args, _ := nr.storage.Builder.
		Select("history_email_notify.request_id", "users.email", "notify_type", "subject", "body").
		From(mailHistoryTable).
		Where("history_email_notify.status = 'processed'").
		InnerJoin("users on history_email_notify.user_id = users.id").
		Limit(batch.BatchSize).
		ToSql()

	rows, err := nr.storage.Pool.Query(ctx, sql, args...)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("%s - nr.storage.Pool.Query: %w", op, postgres.MappingErrors(err))
	}

	mailsData := make([]*dto.MailIdempotencyData, 0)

	for rows.Next() {
		var tmp dto.MailIdempotencyData
		var notifyType int

		err := rows.Scan(&tmp.RequestID, &tmp.Mail, &notifyType, &tmp.Subject, &tmp.Body)
		if err != nil {
			return nil, fmt.Errorf("%s - rows.Scan: %w", op, postgres.MappingErrors(err))
		}
		tmp.NotifyType = models.NotifyType(notifyType)
		mailsData = append(mailsData, &tmp)
	}

	return mailsData, nil
}

func (nr *NotifyRepo) ProcessedBatchMailNotify(ctx context.Context, keys []*dto.IdempotencyKey) error {
	const op = "NotifyRepo.ProcessedBatchMailNotify"
	const spanName = "ProcessedBatchMailNotify"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	stringKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		stringKeys = append(stringKeys, key.RequestID)
	}

	sql, args, _ := nr.storage.Builder.
		Update(mailHistoryTable).
		Set("status", "success").
		Where(squirrel.Eq{"request_id": stringKeys}).
		ToSql()

	_, err := nr.storage.Pool.Exec(ctx, sql, args...)

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - nr.storage.Pool.Exec: %w", op, postgres.MappingErrors(err))
	}

	return nil
}

func (nr *NotifyRepo) GetBatchPhoneNotify(ctx context.Context, batch *dto.BatchNotify) ([]*dto.PhoneIdempotencyData, error) {
	const op = "NotifyRepo.GetBatchPhoneNotify"
	const spanName = "GetBatchPhoneNotify"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	sql, args, _ := nr.storage.Builder.
		Select("history_phone_notify.request_id", "users.phone", "notify_type", "body").
		From(phoneHistoryTable).
		Where("history_phone_notify.status = 'processed'").
		InnerJoin("users on history_phone_notify.user_id = users.id").
		Limit(batch.BatchSize).
		ToSql()

	rows, err := nr.storage.Pool.Query(ctx, sql, args...)

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("%s - nr.storage.Pool.Query: %w", op, postgres.MappingErrors(err))
	}

	phonesData := make([]*dto.PhoneIdempotencyData, 0)

	for rows.Next() {
		var tmp dto.PhoneIdempotencyData
		var notifyType int

		err := rows.Scan(&tmp.RequestID, &tmp.Phone, &notifyType, &tmp.Body)
		if err != nil {
			return nil, fmt.Errorf("%s - rows.Scan: %w", op, postgres.MappingErrors(err))
		}
		tmp.NotifyType = models.NotifyType(notifyType)
		phonesData = append(phonesData, &tmp)
	}

	return phonesData, nil
}

func (nr *NotifyRepo) ProcessedBatchPhoneNotify(ctx context.Context, keys []*dto.IdempotencyKey) error {
	const op = "NotifyRepo.ProcessedBatchPhoneNotify"
	const spanName = "ProcessedBatchPhoneNotify"

	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	stringKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		stringKeys = append(stringKeys, key.RequestID)
	}

	sql, args, _ := nr.storage.Builder.
		Update(phoneHistoryTable).
		Set("status", "success").
		Where(squirrel.Eq{"request_id": stringKeys}).
		ToSql()

	_, err := nr.storage.Pool.Exec(ctx, sql, args...)

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("%s - nr.storage.Pool.Exec: %w", op, postgres.MappingErrors(err))
	}

	return nil
}

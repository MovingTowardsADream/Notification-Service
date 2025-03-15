package notify_test

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"

	"Notification_Service/internal/domain/models"
	repoerr "Notification_Service/internal/infrastructure/repository/errors"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/infrastructure/repository/postgres/notify"
	"Notification_Service/internal/interfaces/dto"
)

func TestNotifyRepoProcessedNotify(t *testing.T) {
	type args struct {
		ctx       context.Context
		processed *dto.ProcessedNotify
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				processed: &dto.ProcessedNotify{
					RequestID: "req_id_b82c4a7d3e1f6a9b8c2d7e4f1a3b6c8",
					UserID:    "9f4e8c1d3a7b6f5e2c8a1b7d3e6f4a9",
					MailDate: &dto.MailDate{
						NotifyType: models.NotifyTypeModerate,
						Subject:    "New alert!",
						Body:       "<html>...",
					},
					PhoneDate: &dto.PhoneDate{
						NotifyType: models.NotifyTypeModerate,
						Body:       "New alert!",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectExec("INSERT INTO history_email_notify").
					WithArgs(
						args.processed.RequestID,
						args.processed.MailDate.NotifyType,
						args.processed.MailDate.Subject,
						args.processed.MailDate.Body,
						"processed",
						args.processed.UserID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectExec("INSERT INTO history_phone_notify").
					WithArgs(
						args.processed.RequestID,
						args.processed.PhoneDate.NotifyType,
						args.processed.PhoneDate.Body,
						"processed",
						args.processed.UserID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "success - only mail",
			args: args{
				ctx: context.Background(),
				processed: &dto.ProcessedNotify{
					RequestID: "req_id_7c3d6f2e8a1a7b3c8d2e6f1a4b9c7",
					UserID:    "6d3a8b2c7e1f4a9d8b3c6a7e2f1d4b",
					MailDate: &dto.MailDate{
						NotifyType: models.NotifyTypeModerate,
						Subject:    "New alert!",
						Body:       "<html>...",
					},
					PhoneDate: nil,
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectExec("INSERT INTO history_email_notify").
					WithArgs(
						args.processed.RequestID,
						args.processed.MailDate.NotifyType,
						args.processed.MailDate.Subject,
						args.processed.MailDate.Body,
						"processed",
						args.processed.UserID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "success - only phone",
			args: args{
				ctx: context.Background(),
				processed: &dto.ProcessedNotify{
					RequestID: "req_id_4e8a1b7c3d6f2e8a1b7c3d6f2e8a1",
					UserID:    "a7b3c8d2e6f1a4b9c7d3e8f2a1b6c",
					MailDate:  nil,
					PhoneDate: &dto.PhoneDate{
						NotifyType: models.NotifyTypeModerate,
						Body:       "New alert!",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectExec("INSERT INTO history_phone_notify").
					WithArgs(
						args.processed.RequestID,
						args.processed.PhoneDate.NotifyType,
						args.processed.PhoneDate.Body,
						"processed",
						args.processed.UserID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "context canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				processed: &dto.ProcessedNotify{
					RequestID: "req_id_72a187b57a357b83216d0018aa47d8c2",
					UserID:    "72a187b57a357b83216d0018aa47d8c2",
					MailDate:  nil,
					PhoneDate: &dto.PhoneDate{
						NotifyType: models.NotifyTypeModerate,
						Body:       "New alert!",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectRollback()
			},
			wantErr: repoerr.ErrCanceled,
		},
		{
			name: "no channels",
			args: args{
				ctx: context.Background(),
				processed: &dto.ProcessedNotify{
					RequestID: "req_id_72a187b57a357b83216d0018aa47d8c2",
					UserID:    "72a187b57a357b83216d0018aa47d8c2",
					MailDate:  nil,
					PhoneDate: nil,
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectCommit()
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()

			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			notifyRepoMock := notify.NewNotifyRepo(postgresMock)

			err := notifyRepoMock.ProcessedNotify(tc.args.ctx, tc.args.processed)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNotifyRepoGetBatchMailNotify(t *testing.T) {
	type args struct {
		ctx   context.Context
		batch *dto.BatchNotify
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         []*dto.MailIdempotencyData
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				batch: &dto.BatchNotify{BatchSize: 1},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"request_id", "email", "notify_type", "subject", "body"}).
					AddRow("req_id_72a187b57a357b83216d0018aa47d8c2", "something@mail.ru", 0, "New alert!", "<html>...</html>")

				m.ExpectQuery(`SELECT history_email_notify.request_id, users.email, notify_type, subject, body` +
					` FROM history_email_notify ` +
					` INNER JOIN users on history_email_notify.user_id = users.id` +
					` WHERE history_email_notify.status = 'processed'` +
					` LIMIT 1`).
					WillReturnRows(rows)
			},
			want: []*dto.MailIdempotencyData{
				&dto.MailIdempotencyData{
					RequestID:  "req_id_72a187b57a357b83216d0018aa47d8c2",
					Mail:       "something@mail.ru",
					NotifyType: 0,
					Subject:    "New alert!",
					Body:       "<html>...</html>",
				},
			},
			wantErr: nil,
		},
		{
			name: "success - not the only one batchSize",
			args: args{
				ctx:   context.Background(),
				batch: &dto.BatchNotify{BatchSize: 3},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"request_id", "email", "notify_type", "subject", "body"}).
					AddRow("req_id_72a187b57a357b83216d0018aa47d8c2", "something@mail.ru", 0, "New alert!", "<html>...</html>").
					AddRow("req_id_5d88ecc17ade4e80b944c91ded7efe83", "gentur@mail.ru", 0, "New alert!", "<html>...</html>")

				m.ExpectQuery(`SELECT history_email_notify.request_id, users.email, notify_type, subject, body` +
					` FROM history_email_notify ` +
					` INNER JOIN users on history_email_notify.user_id = users.id` +
					` WHERE history_email_notify.status = 'processed'` +
					` LIMIT 3`).
					WillReturnRows(rows)
			},
			want: []*dto.MailIdempotencyData{
				&dto.MailIdempotencyData{
					RequestID:  "req_id_72a187b57a357b83216d0018aa47d8c2",
					Mail:       "something@mail.ru",
					NotifyType: 0,
					Subject:    "New alert!",
					Body:       "<html>...</html>",
				},
				&dto.MailIdempotencyData{
					RequestID:  "req_id_5d88ecc17ade4e80b944c91ded7efe83",
					Mail:       "gentur@mail.ru",
					NotifyType: 0,
					Subject:    "New alert!",
					Body:       "<html>...</html>",
				},
			},
			wantErr: nil,
		},
		{
			name: "success - not data",
			args: args{
				ctx:   context.Background(),
				batch: &dto.BatchNotify{BatchSize: 7},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"request_id", "email", "notify_type", "subject", "body"})

				m.ExpectQuery(`SELECT history_email_notify.request_id, users.email, notify_type, subject, body` +
					` FROM history_email_notify ` +
					` INNER JOIN users on history_email_notify.user_id = users.id` +
					` WHERE history_email_notify.status = 'processed'` +
					` LIMIT 7`).
					WillReturnRows(rows)
			},
			want:    []*dto.MailIdempotencyData{},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()

			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			notifyRepoMock := notify.NewNotifyRepo(postgresMock)

			mailData, err := notifyRepoMock.GetBatchMailNotify(tc.args.ctx, tc.args.batch)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, len(tc.want), len(mailData))
			for i := 0; i < len(mailData); i++ {
				assert.Equal(t, *tc.want[i], *mailData[i])
			}

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNotifyRepoGetBatchPhoneNotify(t *testing.T) {
	type args struct {
		ctx   context.Context
		batch *dto.BatchNotify
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         []*dto.PhoneIdempotencyData
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				batch: &dto.BatchNotify{BatchSize: 1},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"request_id", "phone", "notify_type", "body"}).
					AddRow("req_id_72a187b57a357b83216d0018aa47d8c2", "+79135838425", 0, "New alert!")

				m.ExpectQuery(`SELECT history_phone_notify.request_id, users.phone, notify_type, body` +
					` FROM history_phone_notify ` +
					` INNER JOIN users on history_phone_notify.user_id = users.id` +
					` WHERE history_phone_notify.status = 'processed'` +
					` LIMIT 1`).
					WillReturnRows(rows)
			},
			want: []*dto.PhoneIdempotencyData{
				&dto.PhoneIdempotencyData{
					RequestID:  "req_id_72a187b57a357b83216d0018aa47d8c2",
					Phone:      "+79135838425",
					NotifyType: 0,
					Body:       "New alert!",
				},
			},
			wantErr: nil,
		},
		{
			name: "success - not the only one batchSize",
			args: args{
				ctx:   context.Background(),
				batch: &dto.BatchNotify{BatchSize: 3},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"request_id", "phone", "notify_type", "body"}).
					AddRow("req_id_72a187b57a357b83216d0018aa47d8c2", "+79172625406", 0, "New alert!").
					AddRow("req_id_5d88ecc17ade4e80b944c91ded7efe83", "+79958535863", 0, "New alert!")

				m.ExpectQuery(`SELECT history_phone_notify.request_id, users.phone, notify_type, body` +
					` FROM history_phone_notify ` +
					` INNER JOIN users on history_phone_notify.user_id = users.id` +
					` WHERE history_phone_notify.status = 'processed'` +
					` LIMIT 3`).
					WillReturnRows(rows)
			},
			want: []*dto.PhoneIdempotencyData{
				&dto.PhoneIdempotencyData{
					RequestID:  "req_id_72a187b57a357b83216d0018aa47d8c2",
					Phone:      "+79172625406",
					NotifyType: 0,
					Body:       "New alert!",
				},
				&dto.PhoneIdempotencyData{
					RequestID:  "req_id_5d88ecc17ade4e80b944c91ded7efe83",
					Phone:      "+79958535863",
					NotifyType: 0,
					Body:       "New alert!",
				},
			},
			wantErr: nil,
		},
		{
			name: "success - not data",
			args: args{
				ctx:   context.Background(),
				batch: &dto.BatchNotify{BatchSize: 7},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"request_id", "phone", "notify_type", "body"})

				m.ExpectQuery(`SELECT history_phone_notify.request_id, users.phone, notify_type, body` +
					` FROM history_phone_notify ` +
					` INNER JOIN users on history_phone_notify.user_id = users.id` +
					` WHERE history_phone_notify.status = 'processed'` +
					` LIMIT 7`).
					WillReturnRows(rows)
			},
			want:    []*dto.PhoneIdempotencyData{},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()

			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			notifyRepoMock := notify.NewNotifyRepo(postgresMock)

			phoneData, err := notifyRepoMock.GetBatchPhoneNotify(tc.args.ctx, tc.args.batch)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, len(tc.want), len(phoneData))
			for i := 0; i < len(phoneData); i++ {
				assert.Equal(t, *tc.want[i], *phoneData[i])
			}

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNotifyRepoProcessedBatchMailNotify(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []*dto.IdempotencyKey
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				keys: []*dto.IdempotencyKey{
					&dto.IdempotencyKey{
						RequestID: "req_id_72a187b57a357b83216d0018aa47d8c2",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE history_email_notify SET status = \$1 WHERE request_id IN \(\$2\)`).
					WithArgs("success", args.keys[0].RequestID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name: "success - not the only one key",
			args: args{
				ctx: context.Background(),
				keys: []*dto.IdempotencyKey{
					&dto.IdempotencyKey{
						RequestID: "req_id_72a187b57a357b83216d0018aa47d8c2",
					},
					&dto.IdempotencyKey{
						RequestID: "req_id_5d88ecc17ade4e80b944c91ded7efe83",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE history_email_notify SET status = \$1 WHERE request_id IN \(\$2\s*,\s*\$3\)`).
					WithArgs("success", args.keys[0].RequestID, args.keys[1].RequestID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 2))
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			notifyRepoMock := notify.NewNotifyRepo(postgresMock)

			err := notifyRepoMock.ProcessedBatchMailNotify(tc.args.ctx, tc.args.keys)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestNotifyRepoProcessedBatchPhoneNotify(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys []*dto.IdempotencyKey
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				keys: []*dto.IdempotencyKey{
					&dto.IdempotencyKey{
						RequestID: "req_id_72a187b57a357b83216d0018aa47d8c2",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE history_phone_notify SET status = \$1 WHERE request_id IN \(\$2\)`).
					WithArgs("success", args.keys[0].RequestID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name: "success - not the only one key",
			args: args{
				ctx: context.Background(),
				keys: []*dto.IdempotencyKey{
					&dto.IdempotencyKey{
						RequestID: "req_id_72a187b57a357b83216d0018aa47d8c2",
					},
					&dto.IdempotencyKey{
						RequestID: "req_id_5d88ecc17ade4e80b944c91ded7efe83",
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec(`UPDATE history_phone_notify SET status = \$1 WHERE request_id IN \(\$2\s*,\s*\$3\)`).
					WithArgs("success", args.keys[0].RequestID, args.keys[1].RequestID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 2))
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			notifyRepoMock := notify.NewNotifyRepo(postgresMock)

			err := notifyRepoMock.ProcessedBatchPhoneNotify(tc.args.ctx, tc.args.keys)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

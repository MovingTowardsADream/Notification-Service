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

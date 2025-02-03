package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"

	"Notification_Service/internal/domain/models"
	repoerr "Notification_Service/internal/infrastructure/repository/errors"
	"Notification_Service/internal/infrastructure/repository/postgres"
	"Notification_Service/internal/interfaces/dto"
)

func TestUsersRepoGetUserCommunication(t *testing.T) {
	type args struct {
		ctx           context.Context
		communication *dto.IdentificationUserCommunication
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         *dto.UserCommunication
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				communication: &dto.IdentificationUserCommunication{
					ID: "72a187b57a357b83216d0018aa47d8c2",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "email", "phone", "email_notify", "phone_notify"}).
					AddRow(args.communication.ID, "boris_johnson@gmail.com", "+447975556677", true, false)

				m.ExpectQuery(`SELECT users.id, users.email, users.phone, notifications.email_notify, notifications.phone_notify` +
					` FROM users INNER JOIN notifications on users.id = notifications.user_id` +
					` WHERE users.id = \$1`).
					WithArgs(args.communication.ID).
					WillReturnRows(rows)
			},
			want: &dto.UserCommunication{
				ID:        "72a187b57a357b83216d0018aa47d8c2",
				Email:     "boris_johnson@gmail.com",
				Phone:     "+447975556677",
				MailPref:  true,
				PhonePref: false,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				communication: &dto.IdentificationUserCommunication{
					ID: "5d88ecc17ade4e80b944c91ded7efe83",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "email", "phone", "email_notify", "phone_notify"})

				m.ExpectQuery(`SELECT users.id, users.email, users.phone, notifications.email_notify, notifications.phone_notify` +
					` FROM users INNER JOIN notifications on users.id = notifications.user_id` +
					` WHERE users.id = \$1`).
					WithArgs(args.communication.ID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: repoerr.ErrNotFound,
		},
		{
			name: "context canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				communication: &dto.IdentificationUserCommunication{
					ID: "3e96b8c14d2f4912a4e7d059b3af7c68",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery(`SELECT users.id, users.email, users.phone, notifications.email_notify, notifications.phone_notify` +
					` FROM users INNER JOIN notifications on users.id = notifications.user_id` +
					` WHERE users.id = \$1`).
					WithArgs(args.communication.ID).
					WillReturnError(context.Canceled)
			},
			want:    nil,
			wantErr: repoerr.ErrCanceled,
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
			userRepoMock := postgres.NewUsersRepo(postgresMock)

			resp, err := userRepoMock.GetUserCommunication(tc.args.ctx, tc.args.communication)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, *tc.want, *resp)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUsersRepoEditPreferences(t *testing.T) {
	type args struct {
		ctx         context.Context
		preferences *dto.UserPreferences
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
				preferences: &dto.UserPreferences{
					UserID: "72a187b57a357b83216d0018aa47d8c2",
					Preferences: dto.Preferences{
						Mail:  &dto.MailPreference{Approval: true},
						Phone: &dto.PhonePreference{Approval: false},
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"email_notify", "phone_notify"}).
					AddRow(true, true)

				m.ExpectQuery(`SELECT email_notify, phone_notify FROM notifications`).
					WithArgs(args.preferences.UserID).
					WillReturnRows(rows)

				m.ExpectExec(`UPDATE notifications SET email_notify = \$1, phone_notify = \$2 WHERE user_id = \$3`).
					WithArgs(args.preferences.Preferences.Mail.Approval, args.preferences.Preferences.Phone.Approval, args.preferences.UserID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "success phone nil",
			args: args{
				ctx: context.Background(),
				preferences: &dto.UserPreferences{
					UserID: "72a187b57a357b83216d0018aa47d8c2",
					Preferences: dto.Preferences{
						Mail:  nil,
						Phone: &dto.PhonePreference{Approval: true},
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"email_notify", "phone_notify"}).
					AddRow(true, true)

				m.ExpectQuery(`SELECT email_notify, phone_notify FROM notifications`).
					WithArgs(args.preferences.UserID).
					WillReturnRows(rows)

				m.ExpectExec(`UPDATE notifications SET email_notify = \$1, phone_notify = \$2 WHERE user_id = \$3`).
					WithArgs(true, args.preferences.Preferences.Phone.Approval, args.preferences.UserID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "success mail nil",
			args: args{
				ctx: context.Background(),
				preferences: &dto.UserPreferences{
					UserID: "72a187b57a357b83216d0018aa47d8c2",
					Preferences: dto.Preferences{
						Mail:  &dto.MailPreference{Approval: true},
						Phone: nil,
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"email_notify", "phone_notify"}).
					AddRow(true, true)

				m.ExpectQuery(`SELECT email_notify, phone_notify FROM notifications`).
					WithArgs(args.preferences.UserID).
					WillReturnRows(rows)

				m.ExpectExec(`UPDATE notifications SET email_notify = \$1, phone_notify = \$2 WHERE user_id = \$3`).
					WithArgs(args.preferences.Preferences.Mail.Approval, true, args.preferences.UserID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "success all nil",
			args: args{
				ctx: context.Background(),
				preferences: &dto.UserPreferences{
					UserID: "72a187b57a357b83216d0018aa47d8c2",
					Preferences: dto.Preferences{
						Mail:  nil,
						Phone: nil,
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"email_notify", "phone_notify"}).
					AddRow(true, false)

				m.ExpectQuery(`SELECT email_notify, phone_notify FROM notifications`).
					WithArgs(args.preferences.UserID).
					WillReturnRows(rows)

				m.ExpectExec(`UPDATE notifications SET email_notify = \$1, phone_notify = \$2 WHERE user_id = \$3`).
					WithArgs(true, false, args.preferences.UserID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				m.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				ctx: context.Background(),
				preferences: &dto.UserPreferences{
					UserID: "72a187b57a357b83216d0018aa47d8c2",
					Preferences: dto.Preferences{
						Mail:  nil,
						Phone: nil,
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"email_notify", "phone_notify"})

				m.ExpectQuery(`SELECT email_notify, phone_notify FROM notifications`).
					WithArgs(args.preferences.UserID).
					WillReturnRows(rows)

				m.ExpectRollback()
			},
			wantErr: repoerr.ErrNotFound,
		},
		{
			name: "context canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				preferences: &dto.UserPreferences{
					UserID: "72a187b57a357b83216d0018aa47d8c2",
					Preferences: dto.Preferences{
						Mail:  nil,
						Phone: nil,
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectQuery(`SELECT email_notify, phone_notify FROM notifications`).
					WithArgs(args.preferences.UserID).
					WillReturnError(context.Canceled)

				m.ExpectRollback()
			},
			wantErr: repoerr.ErrCanceled,
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
			userRepoMock := postgres.NewUsersRepo(postgresMock)

			err := userRepoMock.EditPreferences(tc.args.ctx, tc.args.preferences)
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

func TestUsersRepoCreate(t *testing.T) {
	const timeStr = "2025-02-02 15:30:45 +03:00"

	createdAt, err := time.Parse("2006-01-02 15:04:05 -07:00", timeStr)
	if err != nil {
		t.Fatalf("error parsing test time: %v", err)
		return
	}

	type args struct {
		ctx      context.Context
		userData *dto.User
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         *models.User
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				userData: &dto.User{
					Username: "petrovich",
					Email:    "petrovich@mail.ru",
					Phone:    "+79035742534",
					Password: "secret",
					Preferences: dto.Preferences{
						Mail: &dto.MailPreference{
							Approval: true,
						},
						Phone: &dto.PhonePreference{
							Approval: true,
						},
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"id", "username", "email", "phone", "time"}).
					AddRow("72a185b57aua0457b83k16dj0187d8c0", args.userData.Username, args.userData.Email, args.userData.Phone, createdAt)

				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.userData.Username, args.userData.Email, args.userData.Phone, args.userData.Password).
					WillReturnRows(rows)

				m.ExpectExec("INSERT INTO notifications").
					WithArgs(args.userData.Preferences.Mail.Approval, args.userData.Preferences.Phone.Approval, "72a185b57aua0457b83k16dj0187d8c0").
					WillReturnResult(pgxmock.NewResult("row add", 1))

				m.ExpectCommit()
			},
			want: &models.User{
				ID:        "72a185b57aua0457b83k16dj0187d8c0",
				Username:  "petrovich",
				Email:     "petrovich@mail.ru",
				Phone:     "+79035742534",
				CreatedAt: createdAt,
			},
			wantErr: nil,
		},
		{
			name: "success incomplete pref",
			args: args{
				ctx: context.Background(),
				userData: &dto.User{
					Username: "cheburek",
					Email:    "cheburek@mail.ru",
					Phone:    "+79560549423",
					Password: "secret",
					Preferences: dto.Preferences{
						Mail: &dto.MailPreference{
							Approval: true,
						},
					},
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"id", "username", "email", "phone", "time"}).
					AddRow("900573k16d47d8b5182a185c7aua0b8j", args.userData.Username, args.userData.Email, args.userData.Phone, createdAt)

				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.userData.Username, args.userData.Email, args.userData.Phone, args.userData.Password).
					WillReturnRows(rows)

				m.ExpectExec("INSERT INTO notifications").
					WithArgs(args.userData.Preferences.Mail.Approval, false, "900573k16d47d8b5182a185c7aua0b8j").
					WillReturnResult(pgxmock.NewResult("row add", 1))

				m.ExpectCommit()
			},
			want: &models.User{
				ID:        "900573k16d47d8b5182a185c7aua0b8j",
				Username:  "cheburek",
				Email:     "cheburek@mail.ru",
				Phone:     "+79560549423",
				CreatedAt: createdAt,
			},
			wantErr: nil,
		},
		{
			name: "success without pref",
			args: args{
				ctx: context.Background(),
				userData: &dto.User{
					Username: "beliash",
					Email:    "beliash@mail.ru",
					Phone:    "+79124052485",
					Password: "secret",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"id", "username", "email", "phone", "time"}).
					AddRow("7d88d413ka782a5c0900b86u7b511a5j", args.userData.Username, args.userData.Email, args.userData.Phone, createdAt)

				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.userData.Username, args.userData.Email, args.userData.Phone, args.userData.Password).
					WillReturnRows(rows)

				m.ExpectExec("INSERT INTO notifications").
					WithArgs(false, false, "7d88d413ka782a5c0900b86u7b511a5j").
					WillReturnResult(pgxmock.NewResult("row add", 1))

				m.ExpectCommit()
			},
			want: &models.User{
				ID:        "7d88d413ka782a5c0900b86u7b511a5j",
				Username:  "beliash",
				Email:     "beliash@mail.ru",
				Phone:     "+79124052485",
				CreatedAt: createdAt,
			},
			wantErr: nil,
		},
		{
			name: "already exists",
			args: args{
				ctx: context.Background(),
				userData: &dto.User{
					Username: "beliash",
					Email:    "beliash@mail.ru",
					Phone:    "+79124052485",
					Password: "secret",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.userData.Username, args.userData.Email, args.userData.Phone, args.userData.Password).
					WillReturnError(&pgconn.PgError{
						Code:    "23505",
						Message: "duplicate key value violates unique constraint",
					})

				m.ExpectCommit()
			},
			want:    nil,
			wantErr: repoerr.ErrAlreadyExists,
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
			userRepoMock := postgres.NewUsersRepo(postgresMock)

			user, err := userRepoMock.Create(tc.args.ctx, tc.args.userData)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, *tc.want, *user)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

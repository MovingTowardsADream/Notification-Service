package tests

import (
	"Notification_Service/internal/entity"
	"Notification_Service/internal/notify/repository/postgresdb"
	"Notification_Service/pkg/postgres"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotifyRepo_GetUserCommunication(t *testing.T) {
	type args struct {
		ctx       context.Context
		id        string
		email     string
		phone     string
		mailPref  bool
		phonePref bool
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCase := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.UserCommunication
		wantErr      error
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				id:        "bhufcxsr6ytui",
				email:     "test@example.com",
				phone:     "1234567890",
				mailPref:  true,
				phonePref: false,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"users.id", "users.email", "users.phone", "notifications.email_notify", "notifications.phone_notify"}).
					AddRow(args.id, args.email, args.phone, args.mailPref, args.phonePref)

				m.ExpectQuery("SELECT users.id, users.email, users.phone, notifications.email_notify, notifications.phone_notify FROM users INNER JOIN notifications on users.id = notifications.user_id").
					WithArgs(args.id).
					WillReturnRows(rows)
			},
			want: entity.UserCommunication{
				ID:        "bhufcxsr6ytui",
				Email:     "test@example.com",
				Phone:     "1234567890",
				MailPref:  true,
				PhonePref: false,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				ctx:       context.Background(),
				id:        "bhufcxsr6ytui",
				email:     "test@example.com",
				phone:     "1234567890",
				mailPref:  true,
				phonePref: false,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT users.id, users.email, users.phone, notifications.email_notify, notifications.phone_notify FROM users INNER JOIN notifications on users.id = notifications.user_id").
					WithArgs(args.id).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.UserCommunication{},
			wantErr: errors.New("user not found"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := postgresdb.NewNotifyRepo(postgresMock)

			got, err := userRepoMock.GetUserCommunication(tc.args.ctx, tc.args.id)
			if tc.wantErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

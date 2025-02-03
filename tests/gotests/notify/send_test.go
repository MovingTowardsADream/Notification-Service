package notify_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/tests/gotests"
)

type SendNotifySuite struct {
	gotests.BaseSuite
}

func (s *SendNotifySuite) SetupTest() {
	s.NewTestContext()
}

func (s *SendNotifySuite) TestSuccess() {
	ctx := context.Background()

	respCreateUser, err := s.Repo.Clients().Users.AddUser(ctx, &notifyv1.AddUserReq{
		Username: "termin",
		Email:    "termin@mail.ru",
		Phone:    "+79248538526",
		Password: "secret_password",
	})

	s.Require().NoError(err)

	respSendNotify, err := s.Repo.Clients().Notify.SendMessage(ctx, &notifyv1.SendMessageReq{
		UserID:     respCreateUser.Id,
		NotifyType: 1,
		Channels: &notifyv1.Channels{
			Mail: &notifyv1.MailNotify{
				Body:    "New alert!",
				Subject: "<html>...</html>",
			},
			Phone: &notifyv1.PhoneNotify{
				Body: "New alert!",
			},
		},
	})

	s.Require().NoError(err)
	s.Require().Equal("success", respSendNotify.Respond)
}

func (s *SendNotifySuite) TestNotFound() {
	ctx := context.Background()

	_, err := s.Repo.Clients().Notify.SendMessage(ctx, &notifyv1.SendMessageReq{
		UserID:     "unknown_id",
		NotifyType: 1,
		Channels: &notifyv1.Channels{
			Mail: &notifyv1.MailNotify{
				Body:    "New alert!",
				Subject: "<html>...</html>",
			},
			Phone: &notifyv1.PhoneNotify{
				Body: "New alert!",
			},
		},
	})

	s.Require().Error(err)
}

func TestUserPrefSuite(t *testing.T) {
	suite.Run(t, &SendNotifySuite{gotests.BaseSuite{Name: t.Name()}})
}

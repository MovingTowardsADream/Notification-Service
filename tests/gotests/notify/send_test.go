package notify_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"

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

	md := metadata.Pairs(
		"X-Request-ID", "req_id_23cebbf2b24914e5f91adb2d77ffc3ah",
	)
	ctxWithMeta := metadata.NewOutgoingContext(ctx, md)

	respSendNotify, err := s.Repo.Clients().Notify.SendMessage(ctxWithMeta, &notifyv1.SendMessageReq{
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

	md := metadata.Pairs(
		"X-Request-ID", "req_id_2df4e7bhbf2b24adbc23c3917f5f91ae",
	)
	ctxWithMeta := metadata.NewOutgoingContext(ctx, md)

	_, err := s.Repo.Clients().Notify.SendMessage(ctxWithMeta, &notifyv1.SendMessageReq{
		UserID:     "unknown_id",
		NotifyType: 1,
		Channels: &notifyv1.Channels{
			Mail: &notifyv1.MailNotify{
				Subject: "New alert!",
				Body:    "<html>...</html>",
			},
			Phone: &notifyv1.PhoneNotify{
				Body: "New alert!",
			},
		},
	})

	s.Require().Error(err)
}

func TestSendNotifySuite(t *testing.T) {
	suite.Run(t, &SendNotifySuite{gotests.BaseSuite{Name: t.Name()}})
}

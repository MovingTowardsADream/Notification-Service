package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/tests/gotests"
)

type EditPreferencesSuite struct {
	gotests.BaseSuite
}

func (s *EditPreferencesSuite) SetupTest() {
	s.NewTestContext()
}

func (s *EditPreferencesSuite) TestSuccess() {
	ctx := context.Background()

	respCreateUser, err := s.Repo.Clients().Users.AddUser(ctx, &notifyv1.AddUserReq{
		Username: "beliash",
		Email:    "beliash@mail.ru",
		Phone:    "+79124052485",
		Password: "secret_password",
	})

	s.Require().NoError(err)

	respEditPref, err := s.Repo.Clients().Users.EditPreferences(ctx, &notifyv1.EditPreferencesReq{
		UserID: respCreateUser.Id,
		Preferences: &notifyv1.Preferences{
			Mail: &notifyv1.MailApproval{
				Approval: true,
			},
			Phone: &notifyv1.PhoneApproval{
				Approval: true,
			},
		},
	})

	s.Require().NoError(err)
	s.Require().Equal("success", respEditPref.Respond)
}

func (s *EditPreferencesSuite) TestNotAllPreferences() {
	ctx := context.Background()

	respCreateUser, err := s.Repo.Clients().Users.AddUser(ctx, &notifyv1.AddUserReq{
		Username: "semen",
		Email:    "semen@mail.ru",
		Phone:    "+79806745634",
		Password: "secret_password",
	})

	s.Require().NoError(err)

	respEditPref, err := s.Repo.Clients().Users.EditPreferences(ctx, &notifyv1.EditPreferencesReq{
		UserID: respCreateUser.Id,
		Preferences: &notifyv1.Preferences{
			Mail: &notifyv1.MailApproval{
				Approval: true,
			},
		},
	})

	s.Require().NoError(err)
	s.Require().Equal("success", respEditPref.Respond)
}

func (s *EditPreferencesSuite) TestNotFound() {
	ctx := context.Background()

	_, err := s.Repo.Clients().Users.EditPreferences(ctx, &notifyv1.EditPreferencesReq{
		UserID: "unknown_id",
		Preferences: &notifyv1.Preferences{
			Mail: &notifyv1.MailApproval{
				Approval: true,
			},
			Phone: &notifyv1.PhoneApproval{
				Approval: true,
			},
		},
	})

	s.Require().Error(err)
	s.Require().Equal(status.Code(err), codes.NotFound)
}

func TestUserPrefSuite(t *testing.T) {
	suite.Run(t, &EditPreferencesSuite{gotests.BaseSuite{Name: t.Name()}})
}

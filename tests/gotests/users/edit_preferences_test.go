package users_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"Notification_Service/tests/gotests"

	notifyv1 "Notification_Service/api/gen/go/notify"
)

type EditPreferencesSuite struct {
	gotests.BaseSuite
}

func (s *EditPreferencesSuite) SetupTest() {
	s.NewTestContext()
}

func (s *EditPreferencesSuite) TestSuccess() {
	ctx := context.Background()

	respEditPref, err := s.Repo.Clients().Users.EditPreferences(ctx, &notifyv1.EditPreferencesReq{
		UserID: "5c036da04d6e74d3b885a2389a2bbb17",
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

	respEditPref, err := s.Repo.Clients().Users.EditPreferences(ctx, &notifyv1.EditPreferencesReq{
		UserID: "5c036da04d6e74d3b885a2389a2bbb17",
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
		UserID: "5c036da04d6e74d3b885a2389a2bbb17",
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
}

func TestUserPrefSuite(t *testing.T) {
	suite.Run(t, &EditPreferencesSuite{gotests.BaseSuite{Name: t.Name()}})
}

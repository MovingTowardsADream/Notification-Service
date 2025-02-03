package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/tests/gotests"
)

type CreateUserSuite struct {
	gotests.BaseSuite
}

func (s *CreateUserSuite) SetupTest() {
	s.NewTestContext()
}

func (s *CreateUserSuite) TestSuccess() {
	ctx := context.Background()

	_, err := s.Repo.Clients().Users.AddUser(ctx, &notifyv1.AddUserReq{
		Username: "sejek",
		Email:    "sejek@mail.ru",
		Phone:    "+79628563752",
		Password: "secret_password",
	})

	s.Require().NoError(err)
}

func TestUserPrefSuite(t *testing.T) {
	suite.Run(t, &CreateUserSuite{gotests.BaseSuite{Name: t.Name()}})
}

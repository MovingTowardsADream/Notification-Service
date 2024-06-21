package tests

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"Notification_Service/tests/suite"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUsersPreferences_GoodOutcome(t *testing.T) {
	ctx, st := suite.New(t)

	response, err := st.UsersClient.UserPreferences(ctx, &notifyv1.UserPreferencesRequest{
		UserId: "c3e72e9a467a8f4d327fyc6ba1c66e7u",
		Preferences: &notifyv1.Preferences{
			Mail: &notifyv1.MailApproval{
				Approval: false,
			},
			Phone: &notifyv1.PhoneApproval{
				Approval: true,
			},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "success", response.Respond)
}

func TestUsersPreferences_NotFound(t *testing.T) {
	ctx, st := suite.New(t)

	_, err := st.UsersClient.UserPreferences(ctx, &notifyv1.UserPreferencesRequest{
		UserId: "id_not_exist",
		Preferences: &notifyv1.Preferences{
			Mail: &notifyv1.MailApproval{
				Approval: true,
			},
			Phone: &notifyv1.PhoneApproval{
				Approval: true,
			},
		},
	})

	require.ErrorContains(t, err, "object not found")
}

// TODO Testing invalid arguments in table driven tests

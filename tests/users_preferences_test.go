package tests

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"Notification_Service/tests/suite"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUsersPreferences(t *testing.T) {
	ctx, st := suite.New(t)

	response, err := st.UsersClient.UserPreferences(ctx, &notifyv1.UserPreferencesRequest{
		UserId: "2a17e4039bfc401e9f5b056f9e3732d2",
		Preferences: &notifyv1.Preferences{
			Mail: &notifyv1.MailApproval{
				Approval: true,
			},
			Phone: &notifyv1.PhoneApproval{
				Approval: true,
			},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "success", response.Respond)
}

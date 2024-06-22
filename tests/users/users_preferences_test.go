package users

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"Notification_Service/tests/users/suite"
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

	require.Error(t, err)
	require.ErrorContains(t, err, "object not found")
}

func TestUsersPreferences_InvalidArgument(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		userId      string
		preferences *notifyv1.Preferences
		expectedErr string
	}{
		{
			name:   "Request with empty UserId",
			userId: "",
			preferences: &notifyv1.Preferences{
				Mail: &notifyv1.MailApproval{
					Approval: true,
				},
				Phone: &notifyv1.PhoneApproval{
					Approval: true,
				},
			},
			expectedErr: "field UserId is required",
		},
		{
			name:        "Request with empty Preferences",
			userId:      "c3e72e9a467a8f4d327fyc6ba1c66e7u",
			preferences: nil,
			expectedErr: "field Preferences is required",
		},
		{
			name:        "Request with empty UserId and Preferences",
			userId:      "",
			preferences: nil,
			expectedErr: "field UserId is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.UsersClient.UserPreferences(ctx, &notifyv1.UserPreferencesRequest{
				UserId:      tt.userId,
				Preferences: tt.preferences,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

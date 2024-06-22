package notify

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	suite2 "Notification_Service/tests/notify/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNotifySending_NotFound(t *testing.T) {
	ctx, st := suite2.New(t)

	_, err := st.SendNotifyClient.SendMessage(ctx, &notifyv1.SendMessageRequest{
		UserId:     "id_not_exist",
		NotifyType: notifyv1.NotifyType_alert,
		Channels: &notifyv1.Channels{
			Mail: &notifyv1.MailNotify{
				Subject: "New alert!",
				Body:    "<html>...",
			},
			Phone: &notifyv1.PhoneNotify{
				Body: "New alert!",
			},
		},
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "object not found")
}

func TestNotifySending_InvalidArgument(t *testing.T) {
	ctx, st := suite2.New(t)

	tests := []struct {
		name        string
		userId      string
		notifyType  notifyv1.NotifyType
		channels    *notifyv1.Channels
		expectedErr string
	}{
		{
			name:       "Request with empty UserId",
			userId:     "",
			notifyType: notifyv1.NotifyType_alert,
			channels: &notifyv1.Channels{
				Mail: &notifyv1.MailNotify{
					Subject: "New alert!",
					Body:    "<html>...",
				},
				Phone: &notifyv1.PhoneNotify{
					Body: "New alert!",
				},
			},
			expectedErr: "field UserId is required",
		},
		{
			name:        "Request with empty Channels",
			userId:      "c3e72e9a467a8f4d327fyc6ba1c66e7u",
			notifyType:  notifyv1.NotifyType_alert,
			channels:    nil,
			expectedErr: "field Channels is required",
		},
		{
			name:       "Request with empty Subject in Mail",
			userId:     "c3e72e9a467a8f4d327fyc6ba1c66e7u",
			notifyType: notifyv1.NotifyType_alert,
			channels: &notifyv1.Channels{
				Mail: &notifyv1.MailNotify{
					Subject: "",
					Body:    "<html>...",
				},
				Phone: &notifyv1.PhoneNotify{
					Body: "New alert!",
				},
			},
			expectedErr: "field Subject is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.SendNotifyClient.SendMessage(ctx, &notifyv1.SendMessageRequest{
				UserId:     tt.userId,
				NotifyType: tt.notifyType,
				Channels:   tt.channels,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

package twilio

import (
	"context"
	"fmt"

	restTwilio "github.com/twilio/twilio-go/rest/api/v2010"

	"Notification_Service/internal/interfaces/dto"
)

type NotifyWorkerUseCase struct {
	PhoneSender *Client
}

func NewNotifyWorker(sender *Client) *NotifyWorkerUseCase {
	return &NotifyWorkerUseCase{PhoneSender: sender}
}

func (n *NotifyWorkerUseCase) SendPhoneSMS(ctx context.Context, notify dto.PhoneDate) error {
	params := &restTwilio.CreateMessageParams{}
	params.SetMessagingServiceSid(n.PhoneSender.ServiceSID)
	params.SetTo(notify.Phone)
	params.SetBody(notify.Body)

	if _, err := n.PhoneSender.RestClient.Api.CreateMessage(params); err != nil {
		return fmt.Errorf("twilio - Client - SendSMS: %w", err)
	}

	return nil
}

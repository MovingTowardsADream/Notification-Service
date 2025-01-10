package twilio

import (
	"context"
	"fmt"

	resttwilio "github.com/twilio/twilio-go/rest/api/v2010"

	"Notification_Service/internal/interfaces/dto"
)

type PhoneSender interface {
	SendPhoneSMS(ctx context.Context, notify dto.PhoneDate) error
}

type NotifyWorkerUseCase struct {
	PhoneSender *Client
}

func NewNotifyWorker(sender *Client) *NotifyWorkerUseCase {
	return &NotifyWorkerUseCase{PhoneSender: sender}
}

func (n *NotifyWorkerUseCase) SendPhoneSMS(ctx context.Context, notify dto.PhoneDate) error {
	params := &resttwilio.CreateMessageParams{}
	params.SetMessagingServiceSid(n.PhoneSender.ServiceSID)
	params.SetTo(notify.Phone)
	params.SetBody(notify.Body)

	if _, err := n.PhoneSender.RestClient.Api.CreateMessage(params); err != nil {
		return fmt.Errorf("twilio - Client - SendSMS: %w", err)
	}

	return nil
}

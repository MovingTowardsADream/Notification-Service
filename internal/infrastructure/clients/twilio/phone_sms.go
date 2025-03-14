package twilio

import (
	"context"
	"fmt"

	resttwilio "github.com/twilio/twilio-go/rest/api/v2010"

	"Notification_Service/internal/interfaces/dto"
)

type WorkerPhone struct {
	sender *Client
}

func NewWorkerPhone(sender *Client) *WorkerPhone {
	return &WorkerPhone{sender: sender}
}

func (n *WorkerPhone) SendPhoneSMS(_ context.Context, notify dto.PhoneDate) error {
	const op = "twilio - SendPhoneSMS"

	params := &resttwilio.CreateMessageParams{}
	params.SetMessagingServiceSid(n.sender.ServiceSID)
	//params.SetTo(notify.Phone)
	params.SetBody(notify.Body)

	if _, err := n.sender.RestClient.Api.CreateMessage(params); err != nil {
		return fmt.Errorf("%s - n.sender.RestClient.Api.CreateMessage: %w", op, err)
	}

	return nil
}

package twilio

import (
	"context"
	"fmt"

	//resttwilio "github.com/twilio/twilio-go/rest/api/v2010"

	"Notification_Service/internal/interfaces/dto"
)

type WorkerPhone struct {
	sender *Client
}

func NewWorkerPhone(sender *Client) *WorkerPhone {
	return &WorkerPhone{sender: sender}
}

func (n *WorkerPhone) SendPhoneSMS(ctx context.Context, notify dto.PhoneDate) error {
	//params := &resttwilio.CreateMessageParams{}
	//params.SetMessagingServiceSid(n.sender.ServiceSID)
	//params.SetTo(notify.Phone)
	//params.SetBody(notify.Body)
	//
	//if _, err := n.sender.RestClient.Api.CreateMessage(params); err != nil {
	//	return fmt.Errorf("twilio - Client - SendSMS: %w", err)
	//}

	fmt.Println("SEND MESSAGE ON PHONE: ", notify)

	return nil
}

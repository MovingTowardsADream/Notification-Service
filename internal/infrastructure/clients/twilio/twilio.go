package twilio

import (
	"github.com/twilio/twilio-go"
)

type Client struct {
	RestClient *twilio.RestClient

	ServiceSID string
}

func NewClient(accountSID, authToken, messagingServiceSID string) *Client {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &Client{
		RestClient: client,
		ServiceSID: messagingServiceSID,
	}
}

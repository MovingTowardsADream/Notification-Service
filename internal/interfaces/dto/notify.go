package dto

import (
	"Notification_Service/internal/domain/models"
)

type MailDate struct {
	Mail       string            `json:"mail"`
	NotifyType models.NotifyType `json:"notify_type"`
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
}

type PhoneDate struct {
	Phone      string            `json:"phone"`
	NotifyType models.NotifyType `json:"notify_type"`
	Body       string            `json:"body"`
}

type MailChannel struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type PhoneChannel struct {
	Body string `json:"body"`
}

type Channels struct {
	Mail  *MailChannel  `json:"email"`
	Phone *PhoneChannel `json:"phone"`
}

type ReqNotification struct {
	UserID     string            `json:"user_id"`
	NotifyType models.NotifyType `json:"notify_type"`
	Channels   Channels          `json:"channels"`
}

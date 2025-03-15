package dto

import (
	"Notification_Service/internal/domain/models"
)

type MailInfo struct {
	RequestID string `json:"request_id"`
	Mail      string `json:"mail"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type PhoneInfo struct {
	RequestID string `json:"request_id"`
	Phone     string `json:"phone"`
	Body      string `json:"body"`
}

type MailDate struct {
	NotifyType models.NotifyType `json:"notify_type"`
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
}

type PhoneDate struct {
	NotifyType models.NotifyType `json:"notify_type"`
	Body       string            `json:"body"`
}

type MailSendData struct {
	Mail       string            `json:"mail"`
	NotifyType models.NotifyType `json:"notify_type"`
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
}

type PhoneSendData struct {
	Phone      string            `json:"phone"`
	NotifyType models.NotifyType `json:"notify_type"`
	Body       string            `json:"body"`
}

type MailIdempotencyData struct {
	RequestID  string            `json:"request_id"`
	Mail       string            `json:"mail"`
	NotifyType models.NotifyType `json:"notify_type"`
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
}

type PhoneIdempotencyData struct {
	RequestID  string            `json:"request_id"`
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
	RequestID  string            `json:"request_id"`
	UserID     string            `json:"user_id"`
	NotifyType models.NotifyType `json:"notify_type"`
	Channels   Channels          `json:"channels"`
}

type ProcessedNotify struct {
	RequestID string     `json:"request_id"`
	UserID    string     `json:"user_id"`
	MailDate  *MailDate  `json:"mail_date"`
	PhoneDate *PhoneDate `json:"phone_date"`
}

type BatchNotify struct {
	BatchSize uint64 `json:"batch_size"`
}

type IdempotencyKey struct {
	RequestID string `json:"request_id"`
}

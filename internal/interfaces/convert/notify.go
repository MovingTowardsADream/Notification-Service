package convert

import (
	"fmt"
	"math"

	"Notification_Service/internal/domain/models"
	"Notification_Service/internal/interfaces/dto"
)

func validateNotifyType(notifyType dto.NotifyType) (models.NotifyType, error) {
	if notifyType < 0 || notifyType > 255 {
		return 0, fmt.Errorf("notifyType value %d is out of uint8 range", notifyType)
	}
	return models.NotifyType(notifyType), nil //nolint:gosec // value checked, overflow is not possible
}

func ToProcessedNotify(notifyRequest *dto.ReqNotification, userCommunication *dto.UserCommunication) *dto.ProcessedNotify {
	procNotify := &dto.ProcessedNotify{
		RequestID: notifyRequest.RequestID,
		UserID:    notifyRequest.UserID,
	}

	if userCommunication.MailPref {
		procNotify.MailDate = toMailDate(notifyRequest)
	}

	if userCommunication.PhonePref {
		procNotify.PhoneDate = toPhoneDate(notifyRequest)
	}

	return procNotify
}

func toMailDate(notifyRequest *dto.ReqNotification) *dto.MailDate {
	return &dto.MailDate{
		NotifyType: notifyRequest.NotifyType,
		Subject:    notifyRequest.Channels.Mail.Subject,
		Body:       notifyRequest.Channels.Mail.Body,
	}
}

func toPhoneDate(notifyRequest *dto.ReqNotification) *dto.PhoneDate {
	return &dto.PhoneDate{
		NotifyType: notifyRequest.NotifyType,
		Body:       notifyRequest.Channels.Phone.Body,
	}
}

func MailIdempotencyDataToMailInfo(mailNotify *dto.MailIdempotencyData) dto.MailInfo {
	return dto.MailInfo{
		RequestID: mailNotify.RequestID,
		Mail:      mailNotify.Mail,
		Subject:   mailNotify.Subject,
		Body:      mailNotify.Body,
	}
}

func PhoneIdempotencyDataToPhoneInfo(phoneNotify *dto.PhoneIdempotencyData) dto.PhoneInfo {
	return dto.PhoneInfo{
		RequestID: phoneNotify.RequestID,
		Phone:     phoneNotify.Phone,
		Body:      phoneNotify.Body,
	}
}

func IntToNotifyType(notifyType int) (models.NotifyType, error) {
	if notifyType < 0 || notifyType > math.MaxUint8 {
		return 0, fmt.Errorf("notifyType %d is out of uint8 range (0-255)", notifyType)
	}
	return models.NotifyType(notifyType), nil
}

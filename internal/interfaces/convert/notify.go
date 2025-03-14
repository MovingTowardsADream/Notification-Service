package convert

import (
	"fmt"

	"Notification_Service/internal/domain/models"
	"Notification_Service/internal/interfaces/dto"
)

func validateNotifyType(notifyType dto.NotifyType) (models.NotifyType, error) {
	if notifyType < 0 || notifyType > 255 {
		return 0, fmt.Errorf("notifyType value %d is out of uint8 range", notifyType)
	}
	return models.NotifyType(notifyType), nil
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

func MailDateToMailInfo(mailNotify *dto.MailDate) dto.MailInfo {
	return dto.MailInfo{
		Subject: mailNotify.Subject,
		Body:    mailNotify.Body,
	}
}

func PhoneDateToPhoneInfo(phoneNotify *dto.PhoneDate) dto.PhoneInfo {
	return dto.PhoneInfo{
		Body: phoneNotify.Body,
	}
}

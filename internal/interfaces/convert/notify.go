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

func ToMailDate(notifyRequest *dto.ReqNotification, userCommunication *dto.UserCommunication) *dto.MailDate {
	return &dto.MailDate{
		Mail:       userCommunication.Email,
		NotifyType: notifyRequest.NotifyType,
		Subject:    notifyRequest.Channels.Mail.Subject,
		Body:       notifyRequest.Channels.Mail.Body,
	}
}

func ToPhoneDate(notifyRequest *dto.ReqNotification, userCommunication *dto.UserCommunication) *dto.PhoneDate {
	return &dto.PhoneDate{
		Phone:      userCommunication.Phone,
		NotifyType: notifyRequest.NotifyType,
		Body:       notifyRequest.Channels.Phone.Body,
	}
}

func MailDateToMailInfo(mailNotify *dto.MailDate) dto.MailInfo {
	return dto.MailInfo{
		Mail:    mailNotify.Mail,
		Subject: mailNotify.Subject,
		Body:    mailNotify.Body,
	}
}

func PhoneDateToPhoneInfo(phoneNotify *dto.PhoneDate) dto.PhoneInfo {
	return dto.PhoneInfo{
		Phone: phoneNotify.Phone,
		Body:  phoneNotify.Body,
	}
}

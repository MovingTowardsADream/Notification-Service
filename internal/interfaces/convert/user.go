package convert

import (
	"Notification_Service/internal/interfaces/dto"
)

func ReqNotifyToIDUserCommunication(notifyRequest *dto.ReqNotification) *dto.IdentificationUserCommunication {
	return &dto.IdentificationUserCommunication{
		ID: notifyRequest.UserID,
	}
}

package convert

import (
	notifyv1 "Notification_Service/api/gen/go/notify"
	"Notification_Service/internal/interfaces/dto"
)

func SendMessageReqToReqNotification(req *notifyv1.SendMessageReq) (*dto.ReqNotification, error) {
	if req == nil {
		return nil, nil
	}

	notifyType, err := validateNotifyType(dto.NotifyType(req.NotifyType))
	if err != nil {
		return nil, err
	}

	return &dto.ReqNotification{
		UserID:     req.UserID,
		NotifyType: notifyType,
		Channels: dto.Channels{
			Mail: &dto.MailChannel{
				Subject: req.GetChannels().GetMail().GetSubject(),
				Body:    req.GetChannels().GetMail().GetBody(),
			},
			Phone: &dto.PhoneChannel{
				Body: req.GetChannels().GetPhone().GetBody(),
			},
		},
	}, nil
}

func EditPreferencesReqToUserPreferences(req *notifyv1.EditPreferencesReq) *dto.UserPreferences {
	if req == nil {
		return nil
	}

	preferences := &dto.UserPreferences{
		UserID: req.UserID,
	}

	if req.Preferences.Mail != nil {
		preferences.Preferences.Mail = &dto.MailPreference{
			Approval: req.GetPreferences().GetMail().GetApproval(),
		}
	}

	if req.Preferences.Phone != nil {
		preferences.Preferences.Phone = &dto.PhonePreference{
			Approval: req.GetPreferences().GetPhone().GetApproval(),
		}
	}

	return preferences
}

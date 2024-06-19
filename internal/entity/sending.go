package entity

type (
	MailChannel struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	PhoneChannel struct {
		Body string `json:"body"`
	}

	Channels struct {
		Mail  MailChannel  `json:"email"`
		Phone PhoneChannel `json:"phone"`
	}

	RequestNotification struct {
		UserId     string   `json:"userId"`
		NotifyType string   `json:"notifyType"`
		Channels   Channels `json:"channels"`
	}
)

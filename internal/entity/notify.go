package entity

type (
	MailDate struct {
		Mail       string `json:"mail"`
		NotifyType string `json:"notify_type"`
		Subject    string `json:"subject"`
		Body       string `json:"body"`
	}

	PhoneDate struct {
		Phone      string `json:"phone"`
		NotifyType string `json:"notify_type"`
		Body       string `json:"body"`
	}
)

package entity

type (
	MailDate struct {
		Mail       string `json:"mail"`
		NotifyType string `json:"notifyType"`
		Subject    string `json:"subject"`
		Body       string `json:"body"`
	}

	PhoneDate struct {
		Phone      string `json:"phone"`
		NotifyType string `json:"notifyType"`
		Body       string `json:"body"`
	}
)

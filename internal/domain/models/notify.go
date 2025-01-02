package models

type Notify struct {
	UserID     string
	NotifyType string
	Channels   struct {
		Email struct {
			Subject string
			Body    string
		}
		Phone struct {
			Body string
		}
	}
}

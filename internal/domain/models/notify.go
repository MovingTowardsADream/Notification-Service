package models

type Notify struct {
	UserID     string
	NotifyType NotifyType
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

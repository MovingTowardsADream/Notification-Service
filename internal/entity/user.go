package entity

type UserCommunication struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	MailPref  bool   `json:"mail_pref"`
	PhonePref bool   `json:"phone_pref"`
}

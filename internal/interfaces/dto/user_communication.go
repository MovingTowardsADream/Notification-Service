package dto

type IdentificationUserCommunication struct {
	ID string `json:"id"`
}

type UserCommunication struct {
	ID        string `json:"id"`
	MailPref  bool   `json:"mail_pref"`
	PhonePref bool   `json:"phone_pref"`
}

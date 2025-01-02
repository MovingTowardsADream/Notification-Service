package dto

type MailPreference struct {
	Approval bool `json:"approval"`
}

type PhonePreference struct {
	Approval bool `json:"approval"`
}

type Preferences struct {
	Mail  *MailPreference  `json:"mail"`
	Phone *PhonePreference `json:"phone"`
}

type UserPreferences struct {
	UserID      string      `json:"user_id"`
	Preferences Preferences `json:"preferences"`
}

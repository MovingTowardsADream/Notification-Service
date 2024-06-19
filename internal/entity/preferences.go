package entity

type (
	MailPreference struct {
		Approval bool `json:"approval"`
	}

	PhonePreference struct {
		Approval bool `json:"approval"`
	}

	Preferences struct {
		Mail  *MailPreference  `json:"mail"`
		Phone *PhonePreference `json:"phone"`
	}

	UserPreferences struct {
		UserId      string      `json:"user_id"`
		Preferences Preferences `json:"preferences"`
	}
)

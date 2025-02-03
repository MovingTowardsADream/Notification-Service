package dto

type User struct {
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	Phone       string      `json:"phone"`
	Password    string      `json:"password"`
	Preferences Preferences `json:"preferences"`
}

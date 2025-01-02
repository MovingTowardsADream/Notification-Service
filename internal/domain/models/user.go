package models

import (
	"time"
)

type User struct {
	ID       string
	Username string
	Email    string
	Phone    string
	Time     time.Time
}

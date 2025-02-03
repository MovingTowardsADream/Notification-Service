package usecase

import (
	"errors"
)

var (
	ErrTimeout       = errors.New("deadline exceeded")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

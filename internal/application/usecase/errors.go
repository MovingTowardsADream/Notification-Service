package usecase

import (
	"context"
	"errors"
)

var (
	ErrTimeout  = context.DeadlineExceeded
	ErrNotFound = errors.New("not found")
)

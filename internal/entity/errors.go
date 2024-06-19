package entity

import (
	"context"
	"errors"
)

var (
	// Request errors
	ErrTimeout  = context.DeadlineExceeded
	ErrNotFound = errors.New("Not found")
)

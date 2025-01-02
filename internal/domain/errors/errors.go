package errors

import (
	"errors"
)

var (
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)

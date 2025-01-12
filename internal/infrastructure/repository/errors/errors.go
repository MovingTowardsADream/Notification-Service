package errors

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrCanceled = errors.New("context canceled")
)

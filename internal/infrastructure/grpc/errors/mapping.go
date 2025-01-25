package errors

import (
	"errors"

	"Notification_Service/internal/application/usecase"
)

func MappingErrors(err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, usecase.ErrTimeout):
		return ErrDeadlineExceeded
	default:
		return ErrInternalServer
	}
}

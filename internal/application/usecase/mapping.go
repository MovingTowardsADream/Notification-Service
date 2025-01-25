package usecase

import (
	"errors"

	repoerr "Notification_Service/internal/infrastructure/repository/errors"
)

func mappingErrors(err error) (bool, error) {
	switch {
	case errors.Is(err, repoerr.ErrNotFound):
		return true, ErrNotFound
	case errors.Is(err, repoerr.ErrCanceled):
		return true, ErrTimeout
	default:
		return false, err
	}
}

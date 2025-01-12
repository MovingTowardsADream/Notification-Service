package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	repoerr "Notification_Service/internal/infrastructure/repository/errors"
)

func mappingErrors(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repoerr.ErrNotFound
	case errors.Is(err, context.Canceled):
		return repoerr.ErrCanceled
	default:
		return err
	}
}

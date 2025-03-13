package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	repoerr "Notification_Service/internal/infrastructure/repository/errors"
)

func MappingErrors(err error) error {
	var pgErr *pgconn.PgError

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repoerr.ErrNotFound
	case errors.Is(err, context.Canceled):
		return repoerr.ErrCanceled
	case errors.As(err, &pgErr):
		return mapPgErrCode(pgErr)
	default:
		return err
	}
}

func mapPgErrCode(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	case "23505":
		return repoerr.ErrAlreadyExists
	default:
		return pgErr
	}
}

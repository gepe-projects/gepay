package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Postgres error codes reference:
// https://www.postgresql.org/docs/current/errcodes-appendix.html

const (
	ErrCodeUniqueViolation     = "23505"
	ErrCodeForeignKeyViolation = "23503"
	ErrCodeNotNullViolation    = "23502"
	ErrCodeExclusionViolation  = "23P01"
)

// IsPgErrCode checks if the given error is a Postgres error with the specified code.
func IsPgErrCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == code
	}
	return false
}

// IsUniqueViolation checks if the error is a unique constraint violation.
func IsUniqueViolation(err error) bool {
	return IsPgErrCode(err, ErrCodeUniqueViolation)
}

// IsForeignKeyViolation checks if the error is a foreign key constraint violation.
func IsForeignKeyViolation(err error) bool {
	return IsPgErrCode(err, ErrCodeForeignKeyViolation)
}

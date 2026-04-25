package pgerr

import (
	"errors"

	"git.appkode.ru/pub/go/failure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	codeUniqueViolation     = "23505"
	codeForeignKeyViolation = "23503"
)

// IsNotFound reports whether err is pgx.ErrNoRows.
func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// IsUniqueViolation reports whether err is a PostgreSQL unique constraint violation.
func IsUniqueViolation(err error) bool {
	return pgErrCode(err) == codeUniqueViolation
}

// IsForeignKeyViolation reports whether err is a PostgreSQL foreign key violation.
func IsForeignKeyViolation(err error) bool {
	return pgErrCode(err) == codeForeignKeyViolation
}

// NotFound converts pgx.ErrNoRows to failure.NotFoundError with the given message.
func NotFound(msg string) error {
	return failure.NewNotFoundError(msg)
}

// UniqueViolation converts a unique constraint error to failure.ConflictError.
func UniqueViolation(msg string) error {
	return failure.NewConflictError(msg)
}

// ForeignKeyViolation converts a foreign key error to failure.InvalidArgumentError.
func ForeignKeyViolation(msg string) error {
	return failure.NewInvalidArgumentError(msg)
}

func pgErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

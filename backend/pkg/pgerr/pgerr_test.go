package pgerr_test

import (
	"errors"
	"fmt"
	"testing"

	"git.appkode.ru/pub/go/failure"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"

	"friday/pkg/pgerr"
)

// -----------------------------------------------------------------------
// IsNotFound
// -----------------------------------------------------------------------

func TestIsNotFound_true(t *testing.T) {
	assert.True(t, pgerr.IsNotFound(pgx.ErrNoRows))
}

func TestIsNotFound_wrappedTrue(t *testing.T) {
	wrapped := fmt.Errorf("wrapped: %w", pgx.ErrNoRows)

	assert.True(t, pgerr.IsNotFound(wrapped))
}

func TestIsNotFound_false(t *testing.T) {
	assert.False(t, pgerr.IsNotFound(errors.New("other error")))
}

// -----------------------------------------------------------------------
// IsUniqueViolation
// -----------------------------------------------------------------------

func TestIsUniqueViolation_true(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &pgconn.PgError{Code: "23505"})

	assert.True(t, pgerr.IsUniqueViolation(err))
}

func TestIsUniqueViolation_false(t *testing.T) {
	assert.False(t, pgerr.IsUniqueViolation(errors.New("other")))
}

func TestIsUniqueViolation_wrongCode(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &pgconn.PgError{Code: "23503"})

	assert.False(t, pgerr.IsUniqueViolation(err))
}

// -----------------------------------------------------------------------
// IsForeignKeyViolation
// -----------------------------------------------------------------------

func TestIsForeignKeyViolation_true(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &pgconn.PgError{Code: "23503"})

	assert.True(t, pgerr.IsForeignKeyViolation(err))
}

func TestIsForeignKeyViolation_false(t *testing.T) {
	assert.False(t, pgerr.IsForeignKeyViolation(errors.New("other")))
}

func TestIsForeignKeyViolation_wrongCode(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &pgconn.PgError{Code: "23505"})

	assert.False(t, pgerr.IsForeignKeyViolation(err))
}

// -----------------------------------------------------------------------
// Error constructors
// -----------------------------------------------------------------------

func TestNotFound_returnsNotFoundError(t *testing.T) {
	err := pgerr.NotFound("record not found")

	assert.True(t, failure.IsNotFoundError(err))
}

func TestUniqueViolation_returnsConflictError(t *testing.T) {
	err := pgerr.UniqueViolation("already exists")

	assert.True(t, failure.IsConflictError(err))
}

func TestForeignKeyViolation_returnsInvalidArgumentError(t *testing.T) {
	err := pgerr.ForeignKeyViolation("referenced id does not exist")

	assert.True(t, failure.IsInvalidArgumentError(err))
}

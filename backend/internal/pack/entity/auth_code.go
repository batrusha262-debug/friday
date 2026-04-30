package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthCode struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Code      string    `db:"code"`
	ExpiresAt time.Time `db:"expires_at"`
	Used      bool      `db:"used"`
	CreatedAt time.Time `db:"created_at"`
}

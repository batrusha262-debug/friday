package entity

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

type Session struct {
	ID        uuid.UUID     `db:"id"`
	UserID    values.UserID `db:"user_id"`
	Token     string        `db:"token"`
	CreatedAt time.Time     `db:"created_at"`
}

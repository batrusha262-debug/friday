package entity

import (
	"time"

	"friday/internal/pack/domain/values"
)

type User struct {
	ID        values.UserID `db:"id"`
	Username  string        `db:"username"`
	CreatedAt time.Time     `db:"created_at"`
}

func (e User) ToDomain() values.User {
	return values.User{
		ID:        e.ID,
		Username:  e.Username,
		CreatedAt: e.CreatedAt,
	}
}

package entity

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

type Pack struct {
	ID        values.PackID `db:"id"`
	Title     string        `db:"title"`
	AuthorID  uuid.UUID     `db:"author_id"`
	CreatedAt time.Time     `db:"created_at"`
}

func (e Pack) ToDomain() values.Pack {
	return values.Pack{
		ID:        e.ID,
		Title:     e.Title,
		AuthorID:  e.AuthorID,
		CreatedAt: e.CreatedAt,
	}
}

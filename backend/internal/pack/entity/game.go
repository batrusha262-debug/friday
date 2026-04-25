package entity

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
)

type Game struct {
	ID              values.GameID       `db:"id"`
	PackID          values.PackID       `db:"pack_id"`
	HostID          uuid.UUID           `db:"host_id"`
	Status          enum.GameStatusEnum `db:"status"`
	CreatedAt       time.Time           `db:"created_at"`
	StartedAt       *time.Time          `db:"started_at"`
	FinishedAt      *time.Time          `db:"finished_at"`
	CurrentPickerID *uuid.UUID          `db:"current_picker_id"`
}

func (e Game) ToDomain() values.Game {
	return values.Game{
		ID:              e.ID,
		PackID:          e.PackID,
		HostID:          e.HostID,
		Status:          e.Status,
		CreatedAt:       e.CreatedAt,
		StartedAt:       e.StartedAt,
		FinishedAt:      e.FinishedAt,
		CurrentPickerID: e.CurrentPickerID,
	}
}

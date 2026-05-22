package entity

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

type MiniGame struct {
	ID             uuid.UUID  `db:"id"`
	GameID         uuid.UUID  `db:"game_id"`
	QuestionID     uuid.UUID  `db:"question_id"`
	ExcludedTeamID *uuid.UUID `db:"excluded_team_id"`
	PosX           int16      `db:"pos_x"`
	PosY           int16      `db:"pos_y"`
	StartedAt      time.Time  `db:"started_at"`
	AppearsAt      time.Time  `db:"appears_at"`
	WinnerTeamID   *uuid.UUID `db:"winner_team_id"`
	FinishedAt     *time.Time `db:"finished_at"`
}

func (e MiniGame) ToDomain() values.MiniGame {
	return values.MiniGame{
		ID:             e.ID,
		GameID:         e.GameID,
		QuestionID:     e.QuestionID,
		ExcludedTeamID: e.ExcludedTeamID,
		PosX:           e.PosX,
		PosY:           e.PosY,
		StartedAt:      e.StartedAt,
		AppearsAt:      e.AppearsAt,
		WinnerTeamID:   e.WinnerTeamID,
		FinishedAt:     e.FinishedAt,
	}
}

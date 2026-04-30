package entity

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

type AnswerClaim struct {
	ID         uuid.UUID  `db:"id"`
	GameID     uuid.UUID  `db:"game_id"`
	QuestionID uuid.UUID  `db:"question_id"`
	TeamID     uuid.UUID  `db:"team_id"`
	ClaimedAt  time.Time  `db:"claimed_at"`
	Status     string     `db:"status"`
	ReviewedAt *time.Time `db:"reviewed_at"`
}

func (e AnswerClaim) ToDomain() values.AnswerClaim {
	return values.AnswerClaim{
		ID:         e.ID,
		GameID:     e.GameID,
		QuestionID: e.QuestionID,
		TeamID:     e.TeamID,
		ClaimedAt:  e.ClaimedAt,
		Status:     e.Status,
		ReviewedAt: e.ReviewedAt,
	}
}

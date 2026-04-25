package entity

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

type GameQuestionState struct {
	ID         uuid.UUID         `db:"id"`
	GameID     values.GameID     `db:"game_id"`
	QuestionID values.QuestionID `db:"question_id"`
	AnsweredBy *uuid.UUID        `db:"answered_by"`
	AnsweredAt *time.Time        `db:"answered_at"`
}

func (e GameQuestionState) ToDomain() values.GameQuestionState {
	var answeredBy *values.GameTeamID

	if e.AnsweredBy != nil {
		id := values.GameTeamID(*e.AnsweredBy)
		answeredBy = &id
	}

	return values.GameQuestionState{
		ID:         e.ID,
		GameID:     e.GameID,
		QuestionID: e.QuestionID,
		AnsweredBy: answeredBy,
		AnsweredAt: e.AnsweredAt,
	}
}

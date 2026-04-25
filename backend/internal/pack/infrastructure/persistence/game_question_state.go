package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) MarkQuestionAnswered(ctx context.Context, gameID, questionID uuid.UUID, answeredBy *uuid.UUID) (entity.GameQuestionState, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_question_states (game_id, question_id, answered_by, answered_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (game_id, question_id) DO UPDATE
		    SET answered_by = EXCLUDED.answered_by,
		        answered_at = EXCLUDED.answered_at
		RETURNING
		    id,
		    game_id,
		    question_id,
		    answered_by,
		    answered_at
		`,
		gameID, questionID, answeredBy,
	)
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("mark question answered: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.GameQuestionState])
	if err != nil {
		if pgerr.IsForeignKeyViolation(err) {
			return entity.GameQuestionState{}, pgerr.ForeignKeyViolation("game, question, or team not found")
		}

		return entity.GameQuestionState{}, fmt.Errorf("mark question answered: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListGameQuestionStates(ctx context.Context, gameID uuid.UUID) ([]entity.GameQuestionState, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    game_id,
		    question_id,
		    answered_by,
		    answered_at
		FROM
		    game_question_states
		WHERE game_id = $1
		`,
		gameID,
	)
	if err != nil {
		return nil, fmt.Errorf("list game question states: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.GameQuestionState])
	if err != nil {
		return nil, fmt.Errorf("list game question states: %w", err)
	}

	return entities, nil
}

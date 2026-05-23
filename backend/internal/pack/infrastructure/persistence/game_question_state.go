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
		    answered_at,
		    revealed_count,
		    timer_started_at,
		    wrong_options
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

func (r *PgRepository) EnsureQuestionState(ctx context.Context, gameID, questionID uuid.UUID) (entity.GameQuestionState, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_question_states (game_id, question_id)
		VALUES ($1, $2)
		ON CONFLICT (game_id, question_id) DO UPDATE
		    SET game_id = EXCLUDED.game_id
		RETURNING
		    id,
		    game_id,
		    question_id,
		    answered_by,
		    answered_at,
		    revealed_count,
		    timer_started_at,
		    wrong_options
		`,
		gameID, questionID,
	)
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("ensure question state: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.GameQuestionState])
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("ensure question state: %w", err)
	}

	return e, nil
}

func (r *PgRepository) RecordWrongOption(ctx context.Context, gameID, questionID uuid.UUID, optionIdx int16) (entity.GameQuestionState, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_question_states (game_id, question_id, wrong_options)
		VALUES ($1, $2, ARRAY[$3]::SMALLINT[])
		ON CONFLICT (game_id, question_id) DO UPDATE
		    SET wrong_options = (
		        SELECT ARRAY(SELECT DISTINCT UNNEST(game_question_states.wrong_options || EXCLUDED.wrong_options))
		    ),
		    timer_started_at = NULL
		RETURNING
		    id,
		    game_id,
		    question_id,
		    answered_by,
		    answered_at,
		    revealed_count,
		    timer_started_at,
		    wrong_options
		`,
		gameID, questionID, optionIdx,
	)
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("record wrong option: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.GameQuestionState])
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("record wrong option: %w", err)
	}

	return e, nil
}

func (r *PgRepository) RevealNextOption(ctx context.Context, gameID, questionID uuid.UUID) (entity.GameQuestionState, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_question_states (game_id, question_id, revealed_count)
		VALUES ($1, $2, 1)
		ON CONFLICT (game_id, question_id) DO UPDATE
		    SET revealed_count = LEAST(game_question_states.revealed_count + 1, 4)
		RETURNING
		    id,
		    game_id,
		    question_id,
		    answered_by,
		    answered_at,
		    revealed_count,
		    timer_started_at,
		    wrong_options
		`,
		gameID, questionID,
	)
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("reveal next option: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.GameQuestionState])
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("reveal next option: %w", err)
	}

	return e, nil
}

func (r *PgRepository) StartQuestionTimer(ctx context.Context, gameID, questionID uuid.UUID) (entity.GameQuestionState, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_question_states (game_id, question_id, timer_started_at)
		VALUES ($1, $2, now())
		ON CONFLICT (game_id, question_id) DO UPDATE
		    SET timer_started_at = now()
		RETURNING
		    id,
		    game_id,
		    question_id,
		    answered_by,
		    answered_at,
		    revealed_count,
		    timer_started_at,
		    wrong_options
		`,
		gameID, questionID,
	)
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("start question timer: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.GameQuestionState])
	if err != nil {
		return entity.GameQuestionState{}, fmt.Errorf("start question timer: %w", err)
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
		    answered_at,
		    revealed_count,
		    timer_started_at,
		    wrong_options
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

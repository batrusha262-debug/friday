package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) CreateGame(ctx context.Context, packID, hostID uuid.UUID) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO games (pack_id, host_id)
		VALUES ($1, $2)
		RETURNING
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		`,
		packID, hostID,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("insert game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsForeignKeyViolation(err) {
			return entity.Game{}, pgerr.ForeignKeyViolation("pack or host not found")
		}

		return entity.Game{}, fmt.Errorf("insert game: %w", err)
	}

	return e, nil
}

func (r *PgRepository) GetGame(ctx context.Context, id uuid.UUID) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		FROM
		    games
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("get game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Game{}, pgerr.NotFound("game not found")
		}

		return entity.Game{}, fmt.Errorf("get game: %w", err)
	}

	return e, nil
}

func (r *PgRepository) FindGameByCode(ctx context.Context, code string) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		FROM
		    games
		WHERE (lower($1) = 'aaaaaaaa' OR lower(id::text) LIKE lower($1) || '%')
		    AND status IN ('waiting', 'active')
		    AND is_open = true
		ORDER BY created_at DESC
		LIMIT 1
		`,
		code,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("find game by code: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Game{}, pgerr.NotFound("game not found")
		}

		return entity.Game{}, fmt.Errorf("find game by code: %w", err)
	}

	return e, nil
}

func (r *PgRepository) FindLatestGameByPack(ctx context.Context, packID uuid.UUID) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		FROM
		    games
		WHERE pack_id = $1
		ORDER BY created_at DESC
		LIMIT 1
		`,
		packID,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("find latest game by pack: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Game{}, pgerr.NotFound("game not found")
		}

		return entity.Game{}, fmt.Errorf("find latest game by pack: %w", err)
	}

	return e, nil
}

func (r *PgRepository) DeleteGame(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM games
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete game: %w", err)
	}

	return nil
}

func (r *PgRepository) StartGame(ctx context.Context, id uuid.UUID) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		WITH random_team AS (
		    SELECT id
		    FROM game_teams
		    WHERE game_id = $1
		    ORDER BY random()
		    LIMIT 1
		)
		UPDATE games
		SET
		    status            = 'active',
		    started_at        = now(),
		    current_picker_id = (SELECT id FROM random_team)
		WHERE id = $1 AND status = 'waiting'
		RETURNING
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		`,
		id,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("start game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Game{}, pgerr.NotFound("game not found or not in waiting status")
		}

		return entity.Game{}, fmt.Errorf("start game: %w", err)
	}

	return e, nil
}

func (r *PgRepository) FinishGame(ctx context.Context, id uuid.UUID) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		UPDATE games
		SET
		    status      = 'finished',
		    finished_at = now()
		WHERE id = $1 AND status = 'active'
		RETURNING
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		`,
		id,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("finish game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Game{}, pgerr.NotFound("game not found or not in active status")
		}

		return entity.Game{}, fmt.Errorf("finish game: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ClaimAnswer(ctx context.Context, gameID, questionID, teamID uuid.UUID) (entity.AnswerClaim, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_answer_claims (game_id, question_id, team_id)
		VALUES ($1, $2, $3)
		RETURNING
		    id,
		    game_id,
		    question_id,
		    team_id,
		    claimed_at,
		    status,
		    reviewed_at
		`,
		gameID, questionID, teamID,
	)
	if err != nil {
		return entity.AnswerClaim{}, fmt.Errorf("insert answer claim: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.AnswerClaim])
	if err != nil {
		return entity.AnswerClaim{}, fmt.Errorf("insert answer claim: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ValidateClaim(ctx context.Context, claimID uuid.UUID, approved bool) (entity.AnswerClaim, error) {
	status := "rejected"
	if approved {
		status = "approved"
	}

	rows, err := r.db.Query(ctx,
		`
		UPDATE game_answer_claims
		SET
		    status      = $2,
		    reviewed_at = now()
		WHERE id = $1
		RETURNING
		    id,
		    game_id,
		    question_id,
		    team_id,
		    claimed_at,
		    status,
		    reviewed_at
		`,
		claimID, status,
	)
	if err != nil {
		return entity.AnswerClaim{}, fmt.Errorf("validate claim: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.AnswerClaim])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.AnswerClaim{}, pgerr.NotFound("claim not found")
		}

		return entity.AnswerClaim{}, fmt.Errorf("validate claim: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListPendingClaims(ctx context.Context, gameID uuid.UUID) ([]entity.AnswerClaim, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    game_id,
		    question_id,
		    team_id,
		    claimed_at,
		    status,
		    reviewed_at
		FROM
		    game_answer_claims
		WHERE
		    game_id = $1
		    AND status = 'pending'
		ORDER BY claimed_at ASC
		`,
		gameID,
	)
	if err != nil {
		return nil, fmt.Errorf("list pending claims: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.AnswerClaim])
	if err != nil {
		return nil, fmt.Errorf("list pending claims: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) SetGameOpen(ctx context.Context, id uuid.UUID, open bool) (entity.Game, error) {
	rows, err := r.db.Query(ctx,
		`
		UPDATE games
		SET is_open = $2
		WHERE id = $1
		RETURNING
		    id,
		    pack_id,
		    host_id,
		    status,
		    is_open,
		    created_at,
		    started_at,
		    finished_at,
		    current_picker_id
		`,
		id, open,
	)
	if err != nil {
		return entity.Game{}, fmt.Errorf("set game open: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Game{}, pgerr.NotFound("game not found")
		}

		return entity.Game{}, fmt.Errorf("set game open: %w", err)
	}

	return e, nil
}

func (r *PgRepository) SetCurrentPicker(ctx context.Context, gameID uuid.UUID, teamID *uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		UPDATE games
		SET current_picker_id = $2
		WHERE id = $1
		`,
		gameID, teamID,
	)
	if err != nil {
		return fmt.Errorf("set current picker: %w", err)
	}

	return nil
}

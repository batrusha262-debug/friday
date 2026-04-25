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

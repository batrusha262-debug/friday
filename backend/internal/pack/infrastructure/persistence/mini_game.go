package persistence

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

const miniGameDelay = 5 * time.Second

func (r *PgRepository) CreateMiniGame(ctx context.Context, gameID, questionID uuid.UUID, excludedTeamID *uuid.UUID) (entity.MiniGame, error) {
	posX := int16(15 + rand.Intn(70))
	posY := int16(20 + rand.Intn(60))
	appearsAt := time.Now().Add(miniGameDelay)

	rows, err := r.db.Query(ctx,
		`
		INSERT INTO mini_games (game_id, question_id, excluded_team_id, pos_x, pos_y, appears_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
		    id,
		    game_id,
		    question_id,
		    excluded_team_id,
		    pos_x,
		    pos_y,
		    started_at,
		    appears_at,
		    winner_team_id,
		    finished_at
		`,
		gameID, questionID, excludedTeamID, posX, posY, appearsAt,
	)
	if err != nil {
		return entity.MiniGame{}, fmt.Errorf("insert mini game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.MiniGame])
	if err != nil {
		return entity.MiniGame{}, fmt.Errorf("insert mini game: %w", err)
	}

	return e, nil
}

func (r *PgRepository) GetActiveMiniGame(ctx context.Context, gameID uuid.UUID) (entity.MiniGame, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    game_id,
		    question_id,
		    excluded_team_id,
		    pos_x,
		    pos_y,
		    started_at,
		    appears_at,
		    winner_team_id,
		    finished_at
		FROM
		    mini_games
		WHERE
		    game_id = $1
		    AND (finished_at IS NULL OR finished_at > now() - interval '3 seconds')
		ORDER BY started_at DESC
		LIMIT 1
		`,
		gameID,
	)
	if err != nil {
		return entity.MiniGame{}, fmt.Errorf("get active mini game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.MiniGame])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.MiniGame{}, pgerr.NotFound("mini game not found")
		}

		return entity.MiniGame{}, fmt.Errorf("get active mini game: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ClaimMiniGame(ctx context.Context, id, teamID uuid.UUID) (entity.MiniGame, error) {
	rows, err := r.db.Query(ctx,
		`
		UPDATE mini_games
		SET
		    winner_team_id = $2,
		    finished_at    = now()
		WHERE
		    id = $1
		    AND finished_at IS NULL
		    AND appears_at <= now()
		    AND (excluded_team_id IS NULL OR excluded_team_id != $2)
		RETURNING
		    id,
		    game_id,
		    question_id,
		    excluded_team_id,
		    pos_x,
		    pos_y,
		    started_at,
		    appears_at,
		    winner_team_id,
		    finished_at
		`,
		id, teamID,
	)
	if err != nil {
		return entity.MiniGame{}, fmt.Errorf("claim mini game: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.MiniGame])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.MiniGame{}, pgerr.NotFound("mini game not claimable")
		}

		return entity.MiniGame{}, fmt.Errorf("claim mini game: %w", err)
	}

	return e, nil
}

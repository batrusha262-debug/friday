package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) AddGameTeam(ctx context.Context, gameID uuid.UUID, name string) (entity.GameTeam, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO game_teams (game_id, name, order_num)
		VALUES (
		    $1, $2,
		    COALESCE((SELECT MAX(order_num) FROM game_teams WHERE game_id = $1), 0) + 1
		)
		RETURNING
		    id,
		    game_id,
		    name,
		    score,
		    order_num
		`,
		gameID, name,
	)
	if err != nil {
		return entity.GameTeam{}, fmt.Errorf("insert game team: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.GameTeam])
	if err != nil {
		if pgerr.IsUniqueViolation(err) {
			return entity.GameTeam{}, pgerr.UniqueViolation("team name already exists in this game")
		}
		if pgerr.IsForeignKeyViolation(err) {
			return entity.GameTeam{}, pgerr.ForeignKeyViolation("game not found")
		}

		return entity.GameTeam{}, fmt.Errorf("insert game team: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListGameTeams(ctx context.Context, gameID uuid.UUID) ([]entity.GameTeam, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    game_id,
		    name,
		    score,
		    order_num
		FROM
		    game_teams
		WHERE game_id = $1
		ORDER BY order_num
		`,
		gameID,
	)
	if err != nil {
		return nil, fmt.Errorf("list game teams: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.GameTeam])
	if err != nil {
		return nil, fmt.Errorf("list game teams: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) RemoveGameTeam(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM game_teams
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("remove game team: %w", err)
	}

	return nil
}

func (r *PgRepository) AwardTeamPoints(ctx context.Context, teamID uuid.UUID, points int) error {
	_, err := r.db.Exec(ctx,
		`
		UPDATE game_teams
		SET score = score + $2
		WHERE id = $1
		`,
		teamID, points,
	)
	if err != nil {
		return fmt.Errorf("award team points: %w", err)
	}

	return nil
}

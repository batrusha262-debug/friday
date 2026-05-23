package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
)

func (r *PgRepository) ListActiveLobbies(ctx context.Context) ([]entity.Lobby, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    g.id          AS game_id,
		    g.pack_id     AS pack_id,
		    p.title       AS pack_title,
		    COUNT(t.id)   AS team_count,
		    g.is_open     AS is_open
		FROM
		    games g
		    JOIN packs p ON p.id = g.pack_id
		    LEFT JOIN game_teams t ON t.game_id = g.id
		WHERE
		    g.status = 'waiting'
		    AND g.is_open = true
		GROUP BY
		    g.id,
		    g.pack_id,
		    p.title,
		    g.is_open
		HAVING COUNT(t.id) > 0
		ORDER BY g.created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list active lobbies: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Lobby])
	if err != nil {
		return nil, fmt.Errorf("list active lobbies: %w", err)
	}

	return entities, nil
}

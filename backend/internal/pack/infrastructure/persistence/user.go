package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) CreateUser(ctx context.Context, username string) (entity.User, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO users (username)
		VALUES ($1)
		RETURNING
		    id,
		    username,
		    created_at
		`,
		username,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("insert user: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		if pgerr.IsUniqueViolation(err) {
			return entity.User{}, pgerr.UniqueViolation("username already taken")
		}

		return entity.User{}, fmt.Errorf("insert user: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListUsers(ctx context.Context) ([]entity.User, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    username,
		    created_at
		FROM
		    users
		ORDER BY created_at ASC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return entities, nil
}

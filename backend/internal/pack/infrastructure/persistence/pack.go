package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) CreatePack(ctx context.Context, title string, authorID uuid.UUID) (entity.Pack, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO packs (title, author_id)
		VALUES ($1, $2)
		RETURNING
		    id,
		    title,
		    author_id,
		    created_at
		`,
		title, authorID,
	)
	if err != nil {
		return entity.Pack{}, fmt.Errorf("insert pack: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Pack])
	if err != nil {
		if pgerr.IsUniqueViolation(err) {
			return entity.Pack{}, pgerr.UniqueViolation("pack already exists")
		}
		if pgerr.IsForeignKeyViolation(err) {
			return entity.Pack{}, pgerr.ForeignKeyViolation("author not found")
		}

		return entity.Pack{}, fmt.Errorf("insert pack: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListPacks(ctx context.Context) ([]entity.Pack, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    title,
		    author_id,
		    created_at
		FROM
		    packs
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list packs: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Pack])
	if err != nil {
		return nil, fmt.Errorf("list packs: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) ListOpenPacks(ctx context.Context) ([]entity.Pack, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    p.id,
		    p.title,
		    p.author_id,
		    p.created_at
		FROM
		    packs p
		WHERE EXISTS (
		    SELECT 1
		    FROM
		        games g
		    WHERE
		        g.pack_id = p.id
		        AND g.is_open = true
		)
		ORDER BY p.created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list open packs: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Pack])
	if err != nil {
		return nil, fmt.Errorf("list open packs: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) GetPack(ctx context.Context, id uuid.UUID) (entity.Pack, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    title,
		    author_id,
		    created_at
		FROM
		    packs
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return entity.Pack{}, fmt.Errorf("get pack: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Pack])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Pack{}, pgerr.NotFound("pack not found")
		}

		return entity.Pack{}, fmt.Errorf("get pack: %w", err)
	}

	return e, nil
}

func (r *PgRepository) DeletePack(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM packs
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete pack: %w", err)
	}

	return nil
}

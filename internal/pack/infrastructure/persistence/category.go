package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) CreateCategory(ctx context.Context, roundID uuid.UUID, name string) (entity.Category, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO categories (round_id, name, order_num)
		VALUES (
			$1, $2,
			COALESCE((SELECT MAX(order_num) FROM categories WHERE round_id = $1), 0) + 1
		)
		RETURNING
		    id,
		    round_id,
		    name,
		    order_num
		`,
		roundID, name,
	)
	if err != nil {
		return entity.Category{}, fmt.Errorf("insert category: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Category])
	if err != nil {
		if pgerr.IsForeignKeyViolation(err) {
			return entity.Category{}, pgerr.ForeignKeyViolation("round not found")
		}

		return entity.Category{}, fmt.Errorf("insert category: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListCategories(ctx context.Context, roundID uuid.UUID) ([]entity.Category, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    round_id,
		    name,
		    order_num
		FROM
		    categories
		WHERE round_id = $1
		ORDER BY order_num
		`,
		roundID,
	)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Category])
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM categories
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}

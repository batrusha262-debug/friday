package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) CreateRound(ctx context.Context, packID uuid.UUID, name string, roundType enum.RoundTypeEnum) (entity.Round, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO rounds (pack_id, name, type, order_num)
		VALUES (
			$1, $2, $3,
			COALESCE((SELECT MAX(order_num) FROM rounds WHERE pack_id = $1), 0) + 1
		)
		RETURNING
		    id,
		    pack_id,
		    name,
		    type,
		    order_num
		`,
		packID, name, roundType,
	)
	if err != nil {
		return entity.Round{}, fmt.Errorf("insert round: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Round])
	if err != nil {
		if pgerr.IsForeignKeyViolation(err) {
			return entity.Round{}, pgerr.ForeignKeyViolation("pack not found")
		}

		return entity.Round{}, fmt.Errorf("insert round: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListRounds(ctx context.Context, packID uuid.UUID) ([]entity.Round, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    pack_id,
		    name,
		    type,
		    order_num
		FROM
		    rounds
		WHERE pack_id = $1
		ORDER BY order_num
		`,
		packID,
	)
	if err != nil {
		return nil, fmt.Errorf("list rounds: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Round])
	if err != nil {
		return nil, fmt.Errorf("list rounds: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) DeleteRound(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM rounds
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete round: %w", err)
	}

	return nil
}

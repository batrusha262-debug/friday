package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/domain/values"
	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) GetQuestion(ctx context.Context, id uuid.UUID) (entity.Question, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    category_id,
		    price,
		    type,
		    question,
		    answer,
		    comment,
		    media_url,
		    order_num
		FROM
		    questions
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return entity.Question{}, fmt.Errorf("get question: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Question])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Question{}, pgerr.NotFound("question not found")
		}

		return entity.Question{}, fmt.Errorf("get question: %w", err)
	}

	return e, nil
}

func (r *PgRepository) CreateQuestion(ctx context.Context, categoryID uuid.UUID, q values.Question) (entity.Question, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO questions (category_id, price, type, question, answer, comment, media_url, order_num)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			COALESCE((SELECT MAX(order_num) FROM questions WHERE category_id = $1), 0) + 1
		)
		RETURNING
		    id,
		    category_id,
		    price,
		    type,
		    question,
		    answer,
		    comment,
		    media_url,
		    order_num
		`,
		categoryID, q.Price, q.Type, q.Question, q.Answer, q.Comment, q.MediaURL,
	)
	if err != nil {
		return entity.Question{}, fmt.Errorf("insert question: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Question])
	if err != nil {
		if pgerr.IsForeignKeyViolation(err) {
			return entity.Question{}, pgerr.ForeignKeyViolation("category not found")
		}

		return entity.Question{}, fmt.Errorf("insert question: %w", err)
	}

	return e, nil
}

func (r *PgRepository) ListQuestions(ctx context.Context, categoryID uuid.UUID) ([]entity.Question, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    id,
		    category_id,
		    price,
		    type,
		    question,
		    answer,
		    comment,
		    media_url,
		    order_num
		FROM
		    questions
		WHERE category_id = $1
		ORDER BY order_num
		`,
		categoryID,
	)
	if err != nil {
		return nil, fmt.Errorf("list questions: %w", err)
	}

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Question])
	if err != nil {
		return nil, fmt.Errorf("list questions: %w", err)
	}

	return entities, nil
}

func (r *PgRepository) UpdateQuestion(ctx context.Context, id uuid.UUID, q values.Question) (entity.Question, error) {
	rows, err := r.db.Query(ctx,
		`
		UPDATE questions
		SET
			price     = $2,
			type      = $3,
			question  = $4,
			answer    = $5,
			comment   = $6,
			media_url = $7
		WHERE id = $1
		RETURNING
		    id,
		    category_id,
		    price,
		    type,
		    question,
		    answer,
		    comment,
		    media_url,
		    order_num
		`,
		id, q.Price, q.Type, q.Question, q.Answer, q.Comment, q.MediaURL,
	)
	if err != nil {
		return entity.Question{}, fmt.Errorf("update question: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Question])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.Question{}, pgerr.NotFound("question not found")
		}

		return entity.Question{}, fmt.Errorf("update question: %w", err)
	}

	return e, nil
}

func (r *PgRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM questions
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}

	return nil
}

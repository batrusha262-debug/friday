package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"friday/internal/pack/entity"
	"friday/pkg/pgerr"
)

func (r *PgRepository) CreateAuthCode(ctx context.Context, email, code string, expiresAt time.Time) (entity.AuthCode, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO auth_codes (email, code, expires_at)
		VALUES ($1, $2, $3)
		RETURNING
		    id,
		    email,
		    code,
		    expires_at,
		    used,
		    created_at
		`,
		email, code, expiresAt,
	)
	if err != nil {
		return entity.AuthCode{}, fmt.Errorf("insert auth code: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.AuthCode])
	if err != nil {
		return entity.AuthCode{}, fmt.Errorf("insert auth code: %w", err)
	}

	return e, nil
}

func (r *PgRepository) UseAuthCode(ctx context.Context, email, code string) (entity.AuthCode, error) {
	rows, err := r.db.Query(ctx,
		`
		UPDATE auth_codes
		SET used = true
		WHERE
		    email      = $1
		    AND code      = $2
		    AND used      = false
		    AND expires_at > now()
		RETURNING
		    id,
		    email,
		    code,
		    expires_at,
		    used,
		    created_at
		`,
		email, code,
	)
	if err != nil {
		return entity.AuthCode{}, fmt.Errorf("use auth code: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.AuthCode])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.AuthCode{}, pgerr.NotFound("invalid or expired code")
		}

		return entity.AuthCode{}, fmt.Errorf("use auth code: %w", err)
	}

	return e, nil
}

func (r *PgRepository) UpsertAdminUser(ctx context.Context, email string) (entity.User, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO users (username, email, role)
		VALUES ($1, $2, 'admin')
		ON CONFLICT (email) WHERE email IS NOT NULL
		DO UPDATE SET role = 'admin'
		RETURNING
		    id,
		    username,
		    email,
		    role,
		    created_at
		`,
		email, email,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("upsert admin user: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return entity.User{}, fmt.Errorf("upsert admin user: %w", err)
	}

	return e, nil
}

func (r *PgRepository) CreateGuestUser(ctx context.Context, username string) (entity.User, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO users (username, role)
		VALUES ($1, 'guest')
		RETURNING
		    id,
		    username,
		    email,
		    role,
		    created_at
		`,
		username,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("insert guest user: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		return entity.User{}, fmt.Errorf("insert guest user: %w", err)
	}

	return e, nil
}

func (r *PgRepository) CreateSession(ctx context.Context, userID uuid.UUID, token string) (entity.Session, error) {
	rows, err := r.db.Query(ctx,
		`
		INSERT INTO sessions (user_id, token)
		VALUES ($1, $2)
		RETURNING
		    id,
		    user_id,
		    token,
		    created_at
		`,
		userID, token,
	)
	if err != nil {
		return entity.Session{}, fmt.Errorf("insert session: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Session])
	if err != nil {
		return entity.Session{}, fmt.Errorf("insert session: %w", err)
	}

	return e, nil
}

func (r *PgRepository) DeleteSession(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx,
		`
		DELETE FROM sessions
		WHERE token = $1
		`,
		token,
	)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (r *PgRepository) GetSessionUser(ctx context.Context, token string) (entity.User, error) {
	rows, err := r.db.Query(ctx,
		`
		SELECT
		    u.id,
		    u.username,
		    u.email,
		    u.role,
		    u.created_at
		FROM
		    sessions s
		    JOIN users u ON u.id = s.user_id
		WHERE s.token = $1
		`,
		token,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("get session user: %w", err)
	}

	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.User])
	if err != nil {
		if pgerr.IsNotFound(err) {
			return entity.User{}, pgerr.NotFound("session not found")
		}

		return entity.User{}, fmt.Errorf("get session user: %w", err)
	}

	return e, nil
}

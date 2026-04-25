package persistence

import "github.com/jackc/pgx/v5/pgxpool"

type PgRepository struct {
	db *pgxpool.Pool
}

func NewPgRepository(db *pgxpool.Pool) *PgRepository {
	return &PgRepository{db: db}
}

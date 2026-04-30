//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"

	"friday/internal/pack/domain/service"
	"friday/internal/pack/infrastructure/persistence"
	packserver "friday/internal/pack/server"
	"friday/pkg/cfg"
)

type Suite struct {
	suite.Suite
	db         *pgxpool.Pool
	srv        *httptest.Server
	repo       *persistence.PgRepository
	http       *Client
	adminToken string
}

func (s *Suite) SetupSuite() {
	ctx := context.Background()
	config := cfg.Load()

	dsn := "host=" + config.Postgres.Host +
		" port=" + config.Postgres.Port +
		" user=" + config.Postgres.User +
		" password=" + config.Postgres.Password +
		" dbname=" + config.Postgres.Database +
		" sslmode=disable"

	db, err := pgxpool.New(ctx, dsn)
	s.Require().NoError(err)
	s.db = db

	sqlDB, err := sql.Open("pgx", dsn)
	s.Require().NoError(err)
	defer sqlDB.Close()

	goose.SetLogger(goose.NopLogger())
	s.Require().NoError(goose.Up(sqlDB, migrationsDir()))

	pgRepo := persistence.NewPgRepository(db)
	h := packserver.NewHandler(service.NewService(pgRepo, nil), nil)
	r := chi.NewRouter()
	h.Register(r)
	s.srv = httptest.NewServer(r)

	s.repo = pgRepo
	s.http = NewClient(s.srv.URL, func() string { return s.adminToken })
}

func (s *Suite) TearDownSuite() {
	s.srv.Close()
	s.db.Close()
}

func (s *Suite) SetupTest() {
	ctx := context.Background()

	_, err := s.db.Exec(ctx,
		`TRUNCATE game_answer_claims, game_question_states, game_teams, games, questions, categories, rounds, packs, auth_codes, users RESTART IDENTITY CASCADE`)
	s.Require().NoError(err)

	s.adminToken = s.seedAdminSession(ctx)
}

func (s *Suite) seedAdminSession(ctx context.Context) string {
	var userID uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO users (username, email, role) VALUES ('admin', 'admin@test.internal', 'admin') RETURNING id`).Scan(&userID)
	s.Require().NoError(err)

	token := uuid.New().String()

	_, err = s.db.Exec(ctx,
		`INSERT INTO sessions (user_id, token) VALUES ($1, $2)`, userID, token)
	s.Require().NoError(err)

	return token
}

func (s *Suite) seedUser(ctx context.Context) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO users (username) VALUES ('test_user') RETURNING id`).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedPack(ctx context.Context, authorID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO packs (title, author_id) VALUES ('Test Pack', $1) RETURNING id`, authorID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedRound(ctx context.Context, packID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO rounds (pack_id, name, type, order_num) VALUES ($1, 'Round 1', 'standard', 1) RETURNING id`, packID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedCategory(ctx context.Context, roundID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO categories (round_id, name, order_num) VALUES ($1, 'Category 1', 1) RETURNING id`, roundID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedQuestion(ctx context.Context, categoryID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO questions (category_id, price, type, question, answer, order_num)
		 VALUES ($1, 100, 'standard', 'Question?', 'Answer', 1) RETURNING id`, categoryID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedGame(ctx context.Context, packID, hostID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO games (pack_id, host_id) VALUES ($1, $2) RETURNING id`, packID, hostID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedStartedGame(ctx context.Context, packID, hostID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO games (pack_id, host_id, status, started_at) VALUES ($1, $2, 'active', now()) RETURNING id`, packID, hostID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func (s *Suite) seedTeam(ctx context.Context, gameID uuid.UUID) uuid.UUID {
	var id uuid.UUID

	err := s.db.QueryRow(ctx,
		`INSERT INTO game_teams (game_id, name, order_num) VALUES ($1, 'Team 1', 1) RETURNING id`, gameID).Scan(&id)
	s.Require().NoError(err)

	return id
}

func migrationsDir() string {
	if dir := os.Getenv("MIGRATIONS_DIR"); dir != "" {
		return dir
	}

	return "../migrations"
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"friday/internal/pack/domain/service"
	"friday/internal/pack/infrastructure/persistence"
	packserver "friday/internal/pack/server"
	"friday/internal/ws"
	"friday/migrations"
	"friday/pkg/cfg"
	"friday/pkg/contextx"
	"friday/pkg/mailer"
	"friday/pkg/postgres"
)

var logger = contextx.DefaultLogger

type Application struct{}

func New() *Application {
	return &Application{}
}

func (app *Application) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	config := cfg.Load()

	pgConfig := postgres.Config{
		Host:     config.Postgres.Host,
		Port:     config.Postgres.Port,
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		Database: config.Postgres.Database,
	}

	db, err := postgres.New(ctx, pgConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	logger.Info("connected to postgres")

	if err := runMigrations(pgConfig.DSN()); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	logger.Info("migrations applied")

	m := mailer.New(
		config.SMTP.Host,
		config.SMTP.Port,
		config.SMTP.Username,
		config.SMTP.Password,
		config.SMTP.From,
	)

	h := packserver.NewHandler(service.NewService(persistence.NewPgRepository(db), m), ws.NewHub())

	r := chi.NewRouter()
	h.Register(r)

	httpServer := &http.Server{
		Addr:    config.HTTPAddr,
		Handler: r,
	}

	logger.Info("starting server", "addr", config.HTTPAddr)

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		return httpServer.Shutdown(context.Background())
	}
}

func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}

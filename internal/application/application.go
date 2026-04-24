package application

import (
	"context"
	"friday/pkg/cfg"
	"friday/pkg/contextx"
	"friday/pkg/postgres"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var logger *slog.Logger = contextx.DefaultLogger

type Application struct {
	db *pgxpool.Pool
}

func New() *Application {
	return &Application{}
}

func (app *Application) Run() error {
	ctx := context.Background()

	config := cfg.Load()

	db, err := postgres.New(ctx, postgres.Config{
		Host:     config.Postgres.Host,
		Port:     config.Postgres.Port,
		User:     config.Postgres.User,
		Password: config.Postgres.Password,
		Database: config.Postgres.Database,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	app.db = db
	logger.Info("connected to postgres")

	return nil
}

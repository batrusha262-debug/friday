package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"friday/internal/pack/infrastructure/persistence"
	packserver "friday/internal/pack/server"
	"friday/internal/pack/domain/service"
	"friday/pkg/cfg"
	"friday/pkg/contextx"
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

	logger.Info("connected to postgres")

	h := packserver.NewHandler(service.NewService(persistence.NewPgRepository(db)))

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

# Project Overview

## Purpose
`friday` — Go backend service with PostgreSQL.

## Tech Stack
- Go 1.26.1
- PostgreSQL via `github.com/jackc/pgx/v5` (pgxpool)
- `log/slog` для логирования
- Docker / Docker Compose

## Structure
```
cmd/main/           — entrypoint (main.go)
internal/application/ — Application struct, Run()
pkg/cfg/            — конфиг через os.Getenv
pkg/contextx/       — WithLogger, LoggerFromContext, EnrichLogger
pkg/postgres/       — New(ctx, Config) *pgxpool.Pool
db/migrations/      — goose SQL-миграции
```

## Environment Variables (.env)
- POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB

include .env
export

# --- Project ---
MODULE        := friday
MIGRATION_DIR := db/migrations
DEFAULT_BRANCH := main
LINT_CONFIG   := .golangci.yml

# --- Tool versions ---
GOOSE_VER       := latest
BETTERALIGN_VER := latest
GOLANGCI_VER    := latest
GOFUMPT_VER     := latest
GCI_VER         := latest
GENUM_VER       := v0.1.3

# --- Goose ---
GOOSE_BUILD_TAGS := no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb
GOOSE_DSN         = host=$(POSTGRES_HOST) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable

.PHONY: build run generate lint format test test-integration test-infrastructure notes \
        goose-create goose-up goose-down install-deps up down logs restart betteralign check

# -----------------------------------------------------------------------
# Docker
# -----------------------------------------------------------------------

## Поднять postgres, применить миграции и запустить приложение
up:
	docker compose up -d postgres
	@echo "Waiting for postgres..."
	@until docker compose exec -T postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do sleep 1; done
	@echo "Postgres is ready"
	$(MAKE) goose-up POSTGRES_HOST=localhost
	docker compose up -d --build app

## Остановить всё и удалить тома
down:
	docker compose down -v

## Перезапустить всё с нуля
restart: down up

## Логи приложения в реальном времени
logs:
	docker compose logs -f app

# -----------------------------------------------------------------------
# Local dev
# -----------------------------------------------------------------------

build:
	go build -o bin/friday ./cmd/main

generate:
	go generate ./...

run:
	go run ./cmd/main

# Единая проверка перед завершением работ
check: format lint test notes

notes:
	@grep -rne "NOTE:" . | grep -v ".idea" | grep -v ".git" || echo "no notes"

lint: betteralign
	golangci-lint run --config $(LINT_CONFIG) --new-from-rev origin/$(DEFAULT_BRANCH)

format:
	gci write -s standard -s default -s "prefix($(MODULE))" --skip-generated .
	gofumpt -l -w -extra .

test:
	go test -race -short ./...
	$(MAKE) goose-up POSTGRES_HOST=localhost
	POSTGRES_HOST=localhost go test -v -tags=integration -race -count=1 ./tests/integration/...
	@echo
	@echo "Integration tests passed"
	

test-infrastructure:
	docker compose up -d postgres
	@echo "Waiting for postgres..."
	@until docker compose exec -T postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do sleep 1; done
	@echo "Postgres is ready"


betteralign:
	betteralign -apply ./... || \
	betteralign -apply ./... || \
	betteralign -apply ./...

# -----------------------------------------------------------------------
# Migrations
# -----------------------------------------------------------------------

# Создание файла миграции. Пример: make goose-create name=init
goose-create:
	$(if $(value name),,$(error Migration name is not specified. Use "make goose-create name=yourname"))
	mkdir -p $(MIGRATION_DIR)
	goose -v -allow-missing -dir $(MIGRATION_DIR) create $(name) sql

goose-up:
	-goose -v -allow-missing -dir $(MIGRATION_DIR) postgres "$(GOOSE_DSN)" up
	goose -v -dir $(MIGRATION_DIR) postgres "$(GOOSE_DSN)" status

goose-down:
	goose -v -dir $(MIGRATION_DIR) postgres "$(GOOSE_DSN)" down
	goose -v -dir $(MIGRATION_DIR) postgres "$(GOOSE_DSN)" status

# -----------------------------------------------------------------------
# Tools
# -----------------------------------------------------------------------

install-deps:
	go install -tags='$(GOOSE_BUILD_TAGS)' github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VER)
	go install github.com/dkorunic/betteralign/cmd/betteralign@$(BETTERALIGN_VER)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VER)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VER)
	go install github.com/daixiang0/gci@$(GCI_VER)
	go install git.appkode.ru/pub/go/genum/cmd/genum@$(GENUM_VER)

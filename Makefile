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

# --- Goose ---
GOOSE_BUILD_TAGS := no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb
GOOSE_DSN        := host=$(POSTGRES_HOST) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable

.PHONY: build run lint format test notes goose-create goose-up goose-down install-deps

build:
	go build -o bin/friday ./cmd/main

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
	go test -race ./...

betteralign:
	betteralign -apply ./... || \
	betteralign -apply ./... || \
	betteralign -apply ./...

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

install-deps:
	go install -tags='$(GOOSE_BUILD_TAGS)' github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VER)
	go install github.com/dkorunic/betteralign/cmd/betteralign@$(BETTERALIGN_VER)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VER)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VER)
	go install github.com/daixiang0/gci@$(GCI_VER)

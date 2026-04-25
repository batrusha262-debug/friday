# Suggested Commands

## Build & Run
```bash
make build          # собрать бинарь в bin/friday
make run            # go run ./cmd/main
docker compose up --build  # запуск с postgres в docker
```

## Code Quality
```bash
make format         # gci + gofumpt
make lint           # betteralign + golangci-lint
make test           # go test -race -short ./...
make test-integration  # интеграционные тесты (нужен postgres + миграции)
make check          # format + lint + test + notes
```

## Migrations (goose)
```bash
make goose-create name=init   # создать миграцию
make goose-up                 # применить миграции
make goose-down               # откатить последнюю
```

## Deps
```bash
make install-deps   # установить все dev-инструменты
go mod tidy
```

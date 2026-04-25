# Style & Conventions

## Code Style
- `gofumpt` + `gci` для форматирования
- `golangci-lint` с конфигом `.golangci.yml`
- `betteralign` для выравнивания полей структур

## Import Order (gci)
1. standard library
2. external packages
3. internal (`prefix(friday)`)

## File Splitting (per entity)
Каждая сущность домена = отдельный файл:
- `handler_{entity}.go` — HTTP-хэндлеры + `{entity}Service` interface
- `service_{entity}.go` — методы Service + `{entity}Repository` interface
- `pg_repository_{entity}.go` — pgx-запросы для сущности
Составные интерфейсы и структуры — в базовом файле (`handler.go`, `service.go`, `pg_repository.go`).

## Error Handling
- Только `git.appkode.ru/pub/go/failure`
- `service_{entity}.go`: `failure.NewInvalidArgumentError`, `failure.NewNotFoundError`
- `pg_repository_{entity}.go`: конвертация через `pkg/pgerr`
- `decode(r, dst)` и `parseID(r, param)` в handler.go → `failure.NewInvalidArgumentError`

## Logging
- `log/slog` через `pkg/contextx`: `WithLogger`, `LoggerFromContext`, `EnrichLogger`
- Default logger: `contextx.DefaultLogger` (TextHandler → stdout)

## Naming
- Короткие имена пакетов: `cfg`, `contextx`, `postgres`, `pgerr`
- Тип ключа контекста — приватная пустая struct: `type contextKeyLogger struct{}`
- Локальные интерфейсы: `{entity}Service`, `{entity}Repository` (unexported)
- Составные интерфейсы: `service`, `repository` (unexported)

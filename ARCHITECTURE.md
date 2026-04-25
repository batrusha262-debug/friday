# Architecture

Проект строится по **Domain-Driven Design** с элементами гексагональной архитектуры.

---

## Структура директорий

```
internal/
  {domain}/                   # один пакет = один bounded context
    model.go                  # сущности и value objects
    repository.go             # публичный составной интерфейс Repository
    service.go                # Service struct + составной repository interface
    pg_repository.go          # PgRepository struct
    handler.go                # Handler struct + Register + составной service interface

    # Файл на каждую сущность:
    service_{entity}.go       # методы Service + локальный {entity}Repository interface
    pg_repository_{entity}.go # методы PgRepository для сущности
    handler_{entity}.go       # HTTP-хэндлеры + локальный {entity}Service interface

  server/
    server.go                 # Server struct, интерфейсы хэндлеров доменов, Handler()

  application/
    application.go            # инициализация всего, graceful shutdown

pkg/
  cfg/        # конфигурация через env
  contextx/   # утилиты контекста (logger, traceID)
  httpx/      # Handler wrapper
    reply/    # HTTP-ответы
  logx/       # slog-атрибуты
  errcodes/   # коды ошибок
  pgerr/      # конвертация pg-ошибок в failure
  postgres/   # инициализация пула соединений
  tests/      # тест-утилиты (Now() и др.)

tests/
  integration/              # интеграционные тесты (build tag: integration)
    suite_test.go           # Suite setup, seed-хелперы
    client_test.go          # HTTP-клиент для тестов
    {entity}_test.go        # тесты эндпоинтов по сущностям
```

---

## Слои внутри домена

### `model.go` — Domain
- Сущности (entity) и value objects
- Доменные перечисления (enum-типы)
- **Без внешних зависимостей** (только `time`, `errors` из stdlib)

### `repository.go` — Port
- Экспортируемый составной интерфейс `Repository`
- Компонуется из per-entity интерфейсов: `Repository interface { packRepository; roundRepository; ... }`

### `service.go` + `service_{entity}.go` — Application Service
- `service.go`: `Service` struct, `NewService`, составной `repository` interface
- `service_{entity}.go`: методы `Service` для сущности + локальный `{entity}Repository` interface
- Валидация через `failure.NewInvalidArgumentError`
- Не знает о HTTP, БД, инфраструктуре

### `pg_repository.go` + `pg_repository_{entity}.go` — Driven Adapter
- `pg_repository.go`: `PgRepository` struct, `NewPgRepository`
- `pg_repository_{entity}.go`: pgx-запросы для сущности
- Конвертация pg-ошибок через `pkg/pgerr`
- Без валидации и бизнес-логики

### `handler.go` + `handler_{entity}.go` — Driving Adapter
- `handler.go`: `Handler` struct, `NewHandler`, `Register`, составной `service` interface, `decode`, `parseID`
- `handler_{entity}.go`: HTTP-хэндлеры + локальный `{entity}Service` interface
- Парсинг запроса → вызов сервиса → `reply.*`

---

## Правила зависимостей

```
handler_{entity}.go  →  {entity}Service interface (определён там же)
handler.go           →  service interface { packService; roundService; ... }
service_{entity}.go  →  {entity}Repository interface (определён там же)
service.go           →  repository interface { packRepo; roundRepo; ... }
pg_repository_{entity}.go →  model.go (только типы) + pkg/pgerr
model.go             →  (ничего внешнего)
```

Слои **не могут** нарушать направление зависимостей:
- `model.go` не импортирует ничего из `internal/`
- `service*.go` не импортирует `pgx`, `chi`, `net/http`
- `handler*.go` не импортирует `pgx`
- `pg_repository*.go` не импортирует `failure`, `chi`

---

## Интерфейсы

Интерфейсы определяются **на стороне потребителя** (Go-идиома).
Каждая сущность — свой интерфейс в своём файле:

```go
// handler_pack.go
type packService interface {
    CreatePack(ctx context.Context, title string, authorID int64) (Pack, error)
    ListPacks(ctx context.Context) ([]Pack, error)
    GetPack(ctx context.Context, id int64) (Pack, error)
    DeletePack(ctx context.Context, id int64) error
}

// handler.go — составной
type service interface {
    packService
    roundService
    categoryService
    questionService
}
```

```go
// service_pack.go
type packRepository interface {
    CreatePack(ctx context.Context, title string, authorID int64) (Pack, error)
    ...
}

// service.go — составной
type repository interface {
    packRepository
    roundRepository
    categoryRepository
    questionRepository
}
```

Публичный `Repository` из `repository.go` — порт домена для внешних потребителей.

---

## Ошибки

- Только `git.appkode.ru/pub/go/failure` — не кастомные типы ошибок
- `service_{entity}.go` создаёт `failure.NewInvalidArgumentError`, `failure.NewNotFoundError`
- `pg_repository_{entity}.go` использует `pkg/pgerr` для конвертации pgx-ошибок
- `handler_{entity}.go` не создаёт ошибки — только пробрасывает через `httpx.Handler`
- `decode` и `parseID` в `handler.go` конвертируют в `failure.NewInvalidArgumentError`

---

## Именование

| Что | Имя |
|-----|-----|
| Публичный интерфейс репозитория | `Repository` |
| Postgres-реализация | `PgRepository` |
| Application service | `Service` |
| HTTP handler | `Handler` |
| Локальный интерфейс в `handler_{entity}.go` | `{entity}Service` (unexported) |
| Составной в `handler.go` | `service` (unexported) |
| Локальный интерфейс в `service_{entity}.go` | `{entity}Repository` (unexported) |
| Составной в `service.go` | `repository` (unexported) |

---

## Server

`internal/server/server.go`:
- По одному **экспортируемому интерфейсу** на домен (`Pack`, `Game`, …)
- `Server` struct с полями этих интерфейсов
- `NewServer(...)` — конструктор в стиле: `func NewServer(pack Pack, ...) Server`
- `Handler() http.Handler` — собирает chi-роутер

`internal/application/application.go` — единственное место где:
- Создаётся `pgxpool.Pool`
- Создаются `PgRepository` → `Service` → `Handler`
- Создаётся `Server` и `http.Server`
- Управляется graceful shutdown через `signal.NotifyContext`

---

## Интеграционные тесты

- Build tag: `//go:build integration`
- Расположение: `tests/integration/`
- Запуск: `make test-integration`
- Один файл на сущность: `{entity}_test.go`
- Структура теста: table-driven, один `testCase` на сценарий
- `Suite.SetupTest()` — TRUNCATE всех таблиц перед каждым тестом
- Seed-хелперы (`seedUser`, `seedPack`, …) — в `suite_test.go`
- HTTP-клиент (`Client`) — в `client_test.go`
- Тесты должны покрывать: success, validation errors (400), not found (404), invalid id (400)

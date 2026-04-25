# Code Style

## SQL-запросы

Каждый SQL-запрос оформляется в многострочном raw-литерале.  
Каждое ключевое слово/клауза — на своей строке, с одним уровнем отступа (табуляция).  
Закрывающий бэктик — на отдельной строке, на уровне аргументов.

Каждая колонка в `SELECT` и `RETURNING` — на своей строке.  
`FROM` — на своей строке, имя таблицы — на следующей строке с отступом.

```go
// ✅ правильно
rows, err := r.db.Query(ctx,
    `
    SELECT
        id,
        title,
        author_id,
        created_at
    FROM
        packs
    WHERE author_id = $1
    ORDER BY created_at DESC
    `,
    authorID,
)

err := r.db.QueryRow(ctx,
    `
    INSERT INTO rounds (pack_id, name, type, order_num)
    VALUES (
        $1, $2, $3,
        COALESCE((SELECT MAX(order_num) FROM rounds WHERE pack_id = $1), 0) + 1
    )
    RETURNING
        id,
        pack_id,
        name,
        type,
        order_num
    `,
    packID, name, roundType,
).Scan(...)

// ❌ неправильно — колонки в одну строку
rows, err := r.db.Query(ctx,
    `
    SELECT id, title, author_id, created_at
    FROM packs
    `,
)

// ❌ неправильно — всё в одну строку
rows, err := r.db.Query(ctx, `SELECT id FROM packs WHERE author_id = $1`, authorID)

// ❌ неправильно — висящий отступ пробелами
rows, err := r.db.Query(ctx,
    `SELECT id, title, author_id, created_at FROM packs
     ORDER BY created_at DESC`)
```

## Пустая строка после объявления переменной

После блока `var` или `:=`-объявлений — пустая строка перед первым использованием.

```go
// ✅ правильно
var p values.Pack

err := r.db.QueryRow(ctx, ...).Scan(&p.ID, &p.Title)

// ✅ правильно
var rounds []values.Round

for rows.Next() {
    var round values.Round

    if err = rows.Scan(...); err != nil {
        return nil, err
    }

    rounds = append(rounds, round)
}

// ❌ неправильно — нет пустой строки
var p values.Pack
err := r.db.QueryRow(ctx, ...).Scan(&p.ID)
```

## Пустая строка перед return

Перед каждым `return` (кроме первого оператора функции или одиночного guard) — пустая строка.

```go
// ✅ правильно
func (r *PgRepository) GetPack(ctx context.Context, id int64) (values.Pack, error) {
    var p values.Pack

    err := r.db.QueryRow(ctx, ...).Scan(&p.ID, &p.Title)
    if err != nil {
        if pgerr.IsNotFound(err) {
            return values.Pack{}, pgerr.NotFound("pack not found")
        }

        return values.Pack{}, fmt.Errorf("get pack: %w", err)
    }

    return p, nil
}

// ✅ guard в начале — пустая строка не нужна
func (s *Service) CreatePack(...) (values.Pack, error) {
    if title == "" {
        return values.Pack{}, failure.NewInvalidArgumentError("title is required")
    }

    return s.repo.CreatePack(ctx, title, authorID)
}

// ❌ неправильно — нет пустой строки перед return
    if err != nil {
        return values.Pack{}, fmt.Errorf("get pack: %w", err)
    }
    return p, nil
```

## Интерфейсы

Параметры — только типы, без имён:

```go
// ✅ правильно
type Repository interface {
    CreatePack(context.Context, string, uuid.UUID) (entity.Pack, error)
    GetPack(context.Context, uuid.UUID) (entity.Pack, error)
    DeletePack(context.Context, uuid.UUID) error
}

// ❌ неправильно — именованные параметры
type Repository interface {
    CreatePack(ctx context.Context, title string, authorID uuid.UUID) (entity.Pack, error)
}
```

## UUID и ID-типы

- Все первичные и внешние ключи в БД — `UUID PRIMARY KEY DEFAULT uuid_generate_v4()` (расширение `uuid-ossp`).
- В `domain/values` для каждой сущности — именованный тип-алиас: `type PackID = uuid.UUID`, `type RoundID = uuid.UUID` и т.д.
- В `entity/` поля ID — `uuid.UUID` с тегом `db:"id"`.
- В интерфейсах Repository/Service для ID-параметров используется `uuid.UUID`.

## Миграции

- Секция `-- +goose Down` **не пишется** — в проекте нет и не будет откатов миграций.
- Пример заголовка DOWN: `-- +goose Down` (пустой или отсутствующий).

## Сводка

| Правило                                                                            | Где применяется                |
| ---------------------------------------------------------------------------------- | ------------------------------ |
| SQL в многострочном литерале, каждая клауза с новой строки                         | `infrastructure/persistence/`  |
| Каждая колонка SELECT/RETURNING — на своей строке; FROM + имя таблицы на следующей | `infrastructure/persistence/`  |
| Параметры интерфейсов — только типы, без имён                                      | `pack.go`, `server/handler.go` |
| ID-типы — `uuid.UUID` в entity/repo/service; алиасы `values.XxxID` в domain        | везде                          |
| Миграции без DROP в секции Down                                                    | `db/migrations/`               |
| Пустая строка после `var x T` и `x := ...`                                         | везде                          |
| Пустая строка перед `return` (если до него есть код)                               | везде                          |

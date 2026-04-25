-- +goose Up
CREATE TYPE round_type AS ENUM ('standard', 'double', 'final');
CREATE TYPE question_type AS ENUM ('standard', 'auction', 'cat_in_bag', 'no_risk');

CREATE TABLE packs (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title      TEXT        NOT NULL,
    author_id  UUID        NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE rounds (
    id        UUID       PRIMARY KEY DEFAULT uuid_generate_v4(),
    pack_id   UUID       NOT NULL REFERENCES packs (id) ON DELETE CASCADE,
    name      TEXT       NOT NULL,
    type      round_type NOT NULL DEFAULT 'standard',
    order_num SMALLINT   NOT NULL,
    UNIQUE (pack_id, order_num)
);

CREATE TABLE categories (
    id        UUID     PRIMARY KEY DEFAULT uuid_generate_v4(),
    round_id  UUID     NOT NULL REFERENCES rounds (id) ON DELETE CASCADE,
    name      TEXT     NOT NULL,
    order_num SMALLINT NOT NULL,
    UNIQUE (round_id, order_num)
);

CREATE TABLE questions (
    id          UUID          PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID          NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    price       INT           NOT NULL,
    type        question_type NOT NULL DEFAULT 'standard',
    question    TEXT          NOT NULL,
    answer      TEXT          NOT NULL,
    comment     TEXT,
    media_url   TEXT,
    order_num   SMALLINT      NOT NULL,
    UNIQUE (category_id, order_num)
);

-- +goose Down
DROP TABLE questions;
DROP TABLE categories;
DROP TABLE rounds;
DROP TABLE packs;
DROP TYPE question_type;
DROP TYPE round_type;

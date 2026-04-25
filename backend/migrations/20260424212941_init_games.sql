-- +goose Up
CREATE TYPE game_status AS ENUM ('waiting', 'active', 'finished');

CREATE TABLE games (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    pack_id     UUID        NOT NULL REFERENCES packs (id),
    host_id     UUID        NOT NULL REFERENCES users (id),
    status      game_status NOT NULL DEFAULT 'waiting',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at  TIMESTAMPTZ,
    finished_at TIMESTAMPTZ
);

CREATE TABLE game_teams (
    id        UUID     PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id   UUID     NOT NULL REFERENCES games (id) ON DELETE CASCADE,
    name      TEXT     NOT NULL,
    score     INT      NOT NULL DEFAULT 0,
    order_num SMALLINT NOT NULL,
    UNIQUE (game_id, name)
);

CREATE TABLE game_question_states (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id     UUID        NOT NULL REFERENCES games (id) ON DELETE CASCADE,
    question_id UUID        NOT NULL REFERENCES questions (id),
    answered_by UUID        REFERENCES game_teams (id),
    answered_at TIMESTAMPTZ,
    UNIQUE (game_id, question_id)
);

-- +goose Down
DROP TABLE game_question_states;
DROP TABLE game_teams;
DROP TABLE games;
DROP TYPE game_status;

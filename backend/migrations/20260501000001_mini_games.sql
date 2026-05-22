-- +goose Up
CREATE TABLE mini_games (
    id               UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id          UUID        NOT NULL REFERENCES games (id) ON DELETE CASCADE,
    question_id      UUID        NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    excluded_team_id UUID                 REFERENCES game_teams (id) ON DELETE SET NULL,
    pos_x            SMALLINT    NOT NULL,
    pos_y            SMALLINT    NOT NULL,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    appears_at       TIMESTAMPTZ NOT NULL,
    winner_team_id   UUID                 REFERENCES game_teams (id) ON DELETE SET NULL,
    finished_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX one_active_mini_game_per_game
    ON mini_games (game_id)
    WHERE finished_at IS NULL;

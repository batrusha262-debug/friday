-- +goose Up
ALTER TABLE game_question_states
    ADD COLUMN revealed_count   SMALLINT    NOT NULL DEFAULT 0,
    ADD COLUMN timer_started_at TIMESTAMPTZ,
    ADD COLUMN wrong_options    SMALLINT[]  NOT NULL DEFAULT ARRAY[]::SMALLINT[];

-- +goose Down

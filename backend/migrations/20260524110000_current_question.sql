-- +goose Up
ALTER TABLE games
    ADD COLUMN current_question_id UUID REFERENCES questions(id) ON DELETE SET NULL;

-- +goose Down

-- +goose Up
ALTER TABLE games ADD COLUMN is_open BOOLEAN NOT NULL DEFAULT false;

-- +goose Down

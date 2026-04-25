-- +goose Up
ALTER TABLE games ADD COLUMN current_picker_id UUID REFERENCES game_teams(id);

-- +goose Down

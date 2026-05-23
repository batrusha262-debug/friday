-- +goose Up
-- +goose StatementBegin

-- Re-seed SpongeBob test game so its join code starts with "67".

DELETE FROM games WHERE id = '66666666-6666-6666-6666-666666666666';

INSERT INTO games (id, pack_id, host_id, status, is_open)
VALUES (
    '67676767-6767-6767-6767-676767676767',
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    'waiting',
    true
)
ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Up

CREATE TABLE game_answer_claims (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id     UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    question_id UUID        NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    team_id     UUID        NOT NULL REFERENCES game_teams(id) ON DELETE CASCADE,
    claimed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    status      TEXT        NOT NULL DEFAULT 'pending',
    reviewed_at TIMESTAMPTZ,
    CONSTRAINT chk_claim_status CHECK (status IN ('pending', 'approved', 'rejected'))
);

-- +goose Down

-- +goose Up
ALTER TABLE questions
    ADD COLUMN options        TEXT[]   NOT NULL DEFAULT ARRAY[]::TEXT[],
    ADD COLUMN correct_option SMALLINT NOT NULL DEFAULT 0;

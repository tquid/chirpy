-- +goose Up
ALTER TABLE users ADD COLUMN is_chirpy_red BOOLEAN NOT NULL default false;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;

-- +goose Up
SELECT 'up SQL query';
ALTER TABLE users
    ADD COLUMN role TEXT NOT NULL DEFAULT 'user'
        CHECK (role IN ('user', 'admin'));
-- +goose Down
SELECT 'down SQL query';
ALTER TABLE users DROP COLUMN role;
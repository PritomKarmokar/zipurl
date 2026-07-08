-- +goose Up
SELECT 'up SQL query';
ALTER TABLE users RENAME COLUMN password_hash TO password;

-- +goose Down
SELECT 'down SQL query';
ALTER TABLE users RENAME COLUMN password TO password_hash;
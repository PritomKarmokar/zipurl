-- +goose Up
SELECT 'up SQL query';
ALTER TABLE urls
    ADD COLUMN max_clicks BIGINT;
-- +goose Down
SELECT 'down SQL query';
ALTER TABLE urls
    DROP COLUMN max_clicks;

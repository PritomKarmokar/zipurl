-- +goose Up
SELECT 'up SQL query';
CREATE TABLE urls (
   id BIGSERIAL PRIMARY KEY,
   url TEXT NOT NULL,
   hashed_token TEXT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
SELECT 'down SQL query';
DROP TABLE urls;

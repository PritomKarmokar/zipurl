-- +goose Up
SELECT 'up SQL query';
ALTER TABLE urls
    ADD COLUMN user_id CHAR(26),
    ADD COLUMN expires_at TIMESTAMPTZ,
    ADD COLUMN click_count BIGINT NOT NULL DEFAULT 0;

ALTER TABLE urls
    ADD CONSTRAINT fk_urls_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE;

CREATE INDEX idx_urls_user ON urls(user_id);
CREATE INDEX idx_urls_expires_at ON urls(expires_at);

-- +goose Down
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_urls_expires_at;
DROP INDEX IF EXISTS idx_urls_user;

ALTER TABLE urls
    DROP CONSTRAINT fk_urls_user;

ALTER TABLE urls
    DROP COLUMN click_count,
    DROP COLUMN expires_at,
    DROP COLUMN user_id;
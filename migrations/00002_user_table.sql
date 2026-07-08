-- +goose Up
SELECT 'up SQL query';
CREATE TABLE users (
    id CHAR(26) PRIMARY KEY,

    username VARCHAR(50) NOT NULL,
    first_name VARCHAR(25) NOT NULL,
    last_name VARCHAR(25) NOT NULL,

    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    status TEXT NOT NULL
       CHECK (status IN ('active', 'inactive', 'blocked'))
                                   DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,
    last_login TIMESTAMPTZ,
    date_joined TIMESTAMPTZ DEFAULT NOW()
);
-- +goose Down
SELECT 'down SQL query';
DROP TABLE users;
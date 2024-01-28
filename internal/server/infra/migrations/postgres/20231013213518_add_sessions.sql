-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id            INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login         TEXT UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    entropy       BYTEA NOT NULL,
    entropy_salt  BYTEA NOT NULL,
    entropy_nonce BYTEA NOT NULL
);
CREATE TABLE sessions (
    id         INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id    INTEGER REFERENCES users(id) NOT NULL,
    token      TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX user_sessions_idx ON sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX user_sessions_idx;
DROP TABLE sessions;
DROP TABLE users;
-- +goose StatementEnd

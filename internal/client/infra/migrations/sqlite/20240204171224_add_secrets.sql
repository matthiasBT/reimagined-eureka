-- +goose Up
-- +goose StatementBegin
-- TODO: create indexes for all tables, including those that are not in this migration
CREATE TABLE credentials (
    id INTEGER PRIMARY KEY,
    purpose TEXT NOT NULL,
    login TEXT NOT NULL,
    encrypted_password BLOB NOT NULL,
    nonce BLOB NOT NULL,
    salt BLOB NOT NULL,
    UNIQUE (purpose, login)
);
CREATE INDEX user_credentials ON credentials(login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

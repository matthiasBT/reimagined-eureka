-- +goose Up
-- +goose StatementBegin
CREATE TABLE credentials (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    login TEXT NOT NULL,
    encrypted_password BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL
);
CREATE TABLE credentials_versions (
    id SERIAL PRIMARY KEY,
    cred_id INTEGER NOT NULL REFERENCES credentials(id),
    version INTEGER NOT NULL,
    meta TEXT NOT NULL,
    login TEXT NOT NULL,
    encrypted_password BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    UNIQUE (cred_id, version)
);

-- TODO: fix all the tables below this line
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    server_id INTEGER,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    UNIQUE (user_id, meta)
);
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    server_id INTEGER,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    UNIQUE (user_id, meta)
);
CREATE TABLE bank_cards (
    id SERIAL PRIMARY KEY,
    server_id INTEGER,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    UNIQUE (user_id, meta)
);
CREATE INDEX user_credentials ON credentials(user_id);
CREATE INDEX user_notes ON notes(user_id);
CREATE INDEX user_files ON files(user_id);
CREATE INDEX user_bank_cards ON bank_cards(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DROP TABLE bank_cards';
SELECT 'DROP TABLE files';
SELECT 'DROP TABLE notes';
SELECT 'DROP TABLE credentials';
-- +goose StatementEnd

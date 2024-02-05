-- +goose Up
-- +goose StatementBegin
-- TODO: create indexes for all tables, including those that are not in this migration
-- TODO: table inheritance?
CREATE TABLE credentials (
    id SERIAL PRIMARY KEY,
    server_id INTEGER UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    login TEXT NOT NULL,
    encrypted_password BLOB NOT NULL,
    nonce BLOB NOT NULL,
    salt BLOB NOT NULL,
    UNIQUE (user_id, server_id)
);
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    server_id INTEGER UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BLOB NOT NULL,
    nonce BLOB NOT NULL,
    salt BLOB NOT NULL,
    UNIQUE (user_id, server_id)
);
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    server_id INTEGER UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BLOB NOT NULL,
    nonce BLOB NOT NULL,
    salt BLOB NOT NULL,
    UNIQUE (user_id, server_id)
);
CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    server_id INTEGER UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BLOB NOT NULL,
    nonce BLOB NOT NULL,
    salt BLOB NOT NULL,
    UNIQUE (user_id, server_id)
);
CREATE INDEX user_credentials ON credentials(user_id);
CREATE INDEX user_notes ON notes(user_id);
CREATE INDEX user_files ON files(user_id);
CREATE INDEX user_cards ON cards(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DROP TABLE cards';
SELECT 'DROP TABLE files';
SELECT 'DROP TABLE notes';
SELECT 'DROP TABLE credentials';
-- +goose StatementEnd

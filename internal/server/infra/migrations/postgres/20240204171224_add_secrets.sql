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
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (cred_id, version)
);
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL
);
CREATE TABLE notes_versions (
    id SERIAL PRIMARY KEY,
    note_id INTEGER NOT NULL REFERENCES notes(id),
    version INTEGER NOT NULL,
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (note_id, version)
);
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL
);
CREATE TABLE files_versions (
    id SERIAL PRIMARY KEY,
    file_id INTEGER NOT NULL REFERENCES files(id),
    version INTEGER NOT NULL,
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (file_id, version)
);
CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL
);
CREATE TABLE cards_versions (
    id SERIAL PRIMARY KEY,
    card_id INTEGER NOT NULL REFERENCES cards(id),
    version INTEGER NOT NULL,
    meta TEXT NOT NULL,
    encrypted_content BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    salt BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (card_id, version)
);
CREATE INDEX user_credentials ON credentials(user_id);
CREATE INDEX user_notes ON notes(user_id);
CREATE INDEX user_files ON files(user_id);
CREATE INDEX user_cards ON cards(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DROP TABLE cards_versions';
SELECT 'DROP TABLE cards';
SELECT 'DROP TABLE files_versions';
SELECT 'DROP TABLE files';
SELECT 'DROP TABLE notes_versions';
SELECT 'DROP TABLE notes';
SELECT 'DROP TABLE credentials_versions';
SELECT 'DROP TABLE credentials';
-- +goose StatementEnd

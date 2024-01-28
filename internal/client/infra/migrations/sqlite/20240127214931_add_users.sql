-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY NOT NULL,
    login TEXT NOT NULL UNIQUE,
    pwd_hash BLOB NOT NULL,
    entropy TEXT NOT NULL,  -- TODO: store hash only
    entropy_encrypted BLOB NOT NULL,
    entropy_salt BLOB NOT NULL,
    entropy_nonce BLOB NOT NULL
);
CREATE TABLE cookies (
    id integer PRIMARY KEY NOT NULL,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id),
    value_encrypted BLOB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cookies;
DROP TABLE master_key_data;
DROP TABLE users;
-- +goose StatementEnd

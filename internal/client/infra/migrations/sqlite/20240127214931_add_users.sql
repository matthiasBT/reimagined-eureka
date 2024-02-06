-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    pwd_hash BLOB NOT NULL,
    entropy_hash BLOB NOT NULL,
    entropy_encrypted BLOB NOT NULL,
    entropy_salt BLOB NOT NULL,
    entropy_nonce BLOB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cookies;
DROP TABLE master_key_data;
DROP TABLE users;
-- +goose StatementEnd

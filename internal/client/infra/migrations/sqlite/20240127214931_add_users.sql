-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY NOT NULL,
    login TEXT NOT NULL UNIQUE CHECK(length(login) >= 6),
    password_hash BLOB NOT NULL,

    secret_cookie_encrypted BLOB,
    master_key_checker_encrypted BLOB NOT NULL,
    master_key_checker TEXT NOT NULL,
    kdf_salt BLOB NOT NULL CHECK(length(kdf_salt) == 16)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

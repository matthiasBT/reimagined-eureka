-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY NOT NULL,
    login TEXT NOT NULL UNIQUE,
    pwd_hash BLOB NOT NULL,
    pwd_salt BLOB NOT NULL

--     kdf_salt BLOB NOT NULL, -- generated for master key and salt kdf
--     master_key_checker TEXT NOT NULL,
--     master_key_checker_encrypted BLOB NOT NULL,
--     secret_cookie_encrypted BLOB  -- TODO: encrypt with master key
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd

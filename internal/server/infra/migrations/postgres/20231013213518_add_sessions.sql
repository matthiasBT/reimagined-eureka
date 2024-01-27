-- +goose Up
-- +goose StatementBegin
create table users(
    id integer primary key generated always as identity,
    login text unique not null,
    password_hash bytea not null
);
create table sessions(
    id integer primary key generated always as identity,
    user_id integer references users(id) not null,
    token text unique not null,
    expires_at timestamptz not null
);
create index user_sessions_idx on sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index user_sessions_idx;
drop table sessions;
drop table users;
-- +goose StatementEnd

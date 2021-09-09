-- +goose Up
CREATE TABLE IF NOT EXISTS conversations (
    id bigserial primary key,
    user_id bigint not null,
    text text not null,
    date timestamp default now()
);
CREATE INDEX IF NOT EXISTS idx_user ON conversations USING hash (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_user;
DROP TABLE IF EXISTS conversations;

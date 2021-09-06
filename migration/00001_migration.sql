-- +goose Up
CREATE TABLE IF NOT EXISTS conversations (
    id bigserial primary key,
    user_id bigint not null,
    text text not null,
    date timestamp default now()
);

-- +goose Down
DROP TABLE conversations;

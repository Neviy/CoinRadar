-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    email VARCHAR(255),
    birthday DATE,
    premium BOOLEAN NOT NULL DEFAULT FALSE,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
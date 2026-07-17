-- +goose Up

CREATE TABLE coins (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(20,8) NOT NULL,
    recommended BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- +goose Down

DROP TABLE coins;
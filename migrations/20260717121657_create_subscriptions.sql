-- +goose Up

CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    coin_id BIGINT NOT NULL,

    interval_minutes INT NOT NULL DEFAULT 10,

    last_sent_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    FOREIGN KEY(coin_id)
        REFERENCES coins(id)
        ON DELETE CASCADE,

    UNIQUE(user_id, coin_id)
);


-- +goose Down

DROP TABLE subscriptions;
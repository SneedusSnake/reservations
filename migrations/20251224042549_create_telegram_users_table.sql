-- +goose Up
CREATE TABLE IF NOT EXISTS telegram_users(
    telegram_id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL
);

-- +goose Down
DROP TABLE telegram_users;

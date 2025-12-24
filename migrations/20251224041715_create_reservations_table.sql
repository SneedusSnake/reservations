-- +goose Up
CREATE TABLE IF NOT EXISTS reservations(
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    start DATETIME NOT NULL,
    end DATETIME NOT NULL
);

-- +goose Down
DROP TABLE reservations;

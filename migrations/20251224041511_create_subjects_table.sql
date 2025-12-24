-- +goose Up
CREATE TABLE IF NOT EXISTS subjects(
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE subjects;

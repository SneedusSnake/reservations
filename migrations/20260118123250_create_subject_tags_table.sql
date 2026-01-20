-- +goose Up
CREATE TABLE IF NOT EXISTS subject_tags (
    subject_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY(subject_id, name)
);

-- +goose Down
DROP TABLE subject_tags;

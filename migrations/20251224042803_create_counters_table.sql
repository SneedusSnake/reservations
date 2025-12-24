-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS counters(
    table_name VARCHAR(255) PRIMARY KEY,
    last_id INTEGER NOT NULL DEFAULT(0)
);

INSERT INTO counters(table_name) VALUES
    ('users'),
    ('subjects'),
    ('reservations')
;
-- +goose StatementEnd

-- +goose Down
DROP TABLE counters;

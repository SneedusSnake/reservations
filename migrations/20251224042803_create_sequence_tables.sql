-- +goose Up
CREATE TABLE IF NOT EXISTS user_seq(
    value INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS subject_seq(
    value INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS reservation_seq(
    value INTEGER PRIMARY KEY
);

INSERT INTO user_seq VALUES (0);

INSERT INTO subject_seq VALUES (0);

INSERT INTO reservation_seq VALUES (0);

-- +goose Down
DROP TABLE user_seq;
DROP TABLE subject_seq;
DROP TABLE reservation_seq;


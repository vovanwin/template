-- +goose Up
CREATE TABLE users (
id serial PRIMARY KEY ,
username TEXT,
name TEXT,
surname TEXT
);

INSERT INTO users VALUES
(0, 'root', '', ''),
(1, 'vojtechvitek', 'Vojtech', 'Vitek');

-- +goose Down
DROP TABLE users;
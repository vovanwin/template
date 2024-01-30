-- +goose Up
CREATE TABLE users
(
    id       serial PRIMARY KEY,
    username TEXT,
    password varchar(255) default ''
);

INSERT INTO users ("username", "password")
VALUES ('root', '$argon2id$v=19$m=65536,t=1,p=12$xri8DaTQcpfqv1MIW1yKxA$z97v+md4xKa+HNs8lKGdD3aDc+H6zQ2ZcUArxUeVTxM'),
       ('vovanwin', '$argon2id$v=19$m=65536,t=1,p=12$xri8DaTQcpfqv1MIW1yKxA$z97v+md4xKa+HNs8lKGdD3aDc+H6zQ2ZcUArxUeVTxM');

-- +goose Down
DROP TABLE users;
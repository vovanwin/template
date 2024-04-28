-- +goose Up
CREATE TABLE users
(
    id         uuid  ,
    username   TEXT,
    password   varchar(255) default '',
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

INSERT INTO users (id, "username", "password")
VALUES ('f296cba6-7d9a-493f-b689-a89b2a2fddf4', 'root',
        '$argon2id$v=19$m=65536,t=1,p=12$xri8DaTQcpfqv1MIW1yKxA$z97v+md4xKa+HNs8lKGdD3aDc+H6zQ2ZcUArxUeVTxM'),
       ('0642aa6b-9da3-4a64-acce-cd4424346b65', 'vovanwin',
        '$argon2id$v=19$m=65536,t=1,p=12$xri8DaTQcpfqv1MIW1yKxA$z97v+md4xKa+HNs8lKGdD3aDc+H6zQ2ZcUArxUeVTxM');

-- +goose Down
DROP TABLE users;
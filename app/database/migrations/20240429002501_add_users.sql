-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users"
(
    "id"         uuid                   NOT NULL primary key,
    "name"       character varying(255) NOT NULL,
    "email"      character varying(255) NOT NULL,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz
);
-- +goose StatementEnd

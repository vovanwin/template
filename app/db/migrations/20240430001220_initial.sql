-- +goose Up
-- Create "users" table
CREATE TABLE "users"
(
    "id"         uuid              NOT NULL,
    "email"      character varying NOT NULL,
    "role"       character varying NOT NULL,
    "password"   character varying NOT NULL,
    "first_name" character varying NOT NULL,
    "last_name"  character varying NOT NULL,
    "deleted_at" timestamptz       NULL,
    "updated_at" timestamptz       NOT NULL default CURRENT_TIMESTAMP,
    "created_at" timestamptz       NOT NULL default CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");

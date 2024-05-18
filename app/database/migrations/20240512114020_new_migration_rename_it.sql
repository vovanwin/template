-- Modify "users" table
ALTER TABLE "users"
    ALTER COLUMN "updated_at" SET NOT NULL,
    ADD COLUMN "login"    character varying NOT NULL,
    ADD COLUMN "password" character varying NOT NULL;

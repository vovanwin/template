-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "age" bigint NULL;
-- Create "posts" table
CREATE TABLE "posts" ("id" uuid NOT NULL, "test" character varying NOT NULL, "deleted_at" timestamptz NULL, "updated_at" timestamptz NULL, "created_at" timestamptz NOT NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "posts_users_posts" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);

-- +goose Up
ALTER TABLE reminders ADD COLUMN require_confirmation BOOLEAN DEFAULT false;
ALTER TABLE reminders ADD COLUMN repeat_interval_minutes INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE reminders DROP COLUMN IF EXISTS repeat_interval_minutes;
ALTER TABLE reminders DROP COLUMN IF EXISTS require_confirmation;

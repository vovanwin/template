-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions (
    token VARCHAR(255) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Индекс для автоматической очистки истекших сессий
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions(expiry);

-- Функция для автоматической очистки истекших сессий
CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions WHERE expiry < NOW();
END;
$$ LANGUAGE plpgsql;

-- Создание расширения pg_cron если доступно (опционально)
-- CREATE EXTENSION IF NOT EXISTS pg_cron;
-- SELECT cron.schedule('cleanup-sessions', '0 */6 * * *', 'SELECT cleanup_expired_sessions();');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT cron.unschedule('cleanup-sessions');
DROP FUNCTION IF EXISTS cleanup_expired_sessions();
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd

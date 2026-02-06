-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) DEFAULT 'user',
    tenant_id VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    settings JSONB DEFAULT '{}',
    components JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Вставляем тестовых пользователей (пароли: password и 123456)
INSERT INTO users (email, password_hash, first_name, last_name, role) VALUES
('admin@example.com', '$argon2id$v=19$m=65536,t=1,p=12$e7hiycNCJOQCI9TH8oJOMQ$BpmBj/C/W/SJJ9feV8o1/WDLjPu5hY3C6/jD6pVZbOQ', 'Admin', 'User', 'admin'),
('user@example.com', '$argon2id$v=19$m=65536,t=1,p=12$E8pTuGmXMt+XOGP63gAbEA$ayMMZIeTbzgxtAYlPCy58yguBzlixxIIlmrtMcpUygw', 'Regular', 'User', 'user')
ON CONFLICT (email) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

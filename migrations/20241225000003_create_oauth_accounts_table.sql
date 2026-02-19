-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS oauth_accounts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider   VARCHAR(50)  NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    email      VARCHAR(255),
    created_at TIMESTAMPTZ  DEFAULT NOW(),
    UNIQUE (provider, provider_id)
);

CREATE INDEX IF NOT EXISTS idx_oauth_accounts_user_id ON oauth_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_accounts_provider ON oauth_accounts(provider, provider_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth_accounts;
-- +goose StatementEnd

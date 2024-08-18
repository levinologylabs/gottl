-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_identity_providers (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider_name VARCHAR(50) NOT NULL CHECK (provider_name IN ('google')),
    provider_user_id VARCHAR(254) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_user_identity_providers_provider_user_id_and_provider_name ON user_identity_providers (provider_user_id, provider_name);

CREATE TABLE IF NOT EXISTS user_identity_provider_state (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    token BYTEA NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_user_identity_provider_state_token ON user_identity_provider_state(token);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_identity_providers;

DROP TABLE IF EXISTS user_identity_provider_state;

-- +goose StatementEnd

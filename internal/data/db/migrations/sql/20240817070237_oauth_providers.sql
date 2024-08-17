-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_identity_providers (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider_name VARCHAR(50) NOT NULL CHECK (provider_name IN ('google')),
    provider_user_id VARCHAR(254) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_user_identity_providers_provider_user_id_and_provider_name ON user_identity_providers (provider_user_id, provider_name);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_identity_providers;

-- +goose StatementEnd

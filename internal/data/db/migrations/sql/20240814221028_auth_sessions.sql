-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_sessions (
    -- table column defaults
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- user identity
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    -- session information
    token BYTEA NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS user_sessions_token_idx ON user_sessions (token);

CREATE INDEX IF NOT EXISTS user_sessions_expires_at_idx ON user_sessions (expires_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_sessions;

-- +goose StatementEnd

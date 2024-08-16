-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_action_tokens (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token BYTEA NOT NULL,
    action VARCHAR(255) NOT NULL CHECK (TYPE IN ('password_reset'))
);

CREATE UNIQUE INDEX IF NOT EXISTS user_action_tokens_token_idx ON user_action_tokens (token);

CREATE INDEX IF NOT EXISTS user_action_tokens_expires_at_idx ON user_action_tokens (expires_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_action_tokens;

-- +goose StatementEnd

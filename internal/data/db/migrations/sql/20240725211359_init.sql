-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    -- table column defaults
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- user identity
    username VARCHAR(100) NOT NULL DEFAULT '',
    email VARCHAR(254) NOT NULL UNIQUE,
    password_hash VARCHAR(500) NOT NULL,
    -- billing information
    stripe_customer_id VARCHAR(50),
    stripe_subscription_id VARCHAR(50),
    subscription_start_date TIMESTAMP,
    subscription_ended_date TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd

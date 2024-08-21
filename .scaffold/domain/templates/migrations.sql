-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS {{ .Computed.sql_table }} (
    -- table column defaults
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS {{ .Computed.sql_table }};

-- +goose StatementEnd

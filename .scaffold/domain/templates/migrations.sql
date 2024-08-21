-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS {{ .Scaffold.sql_table }} (
    -- table column defaults
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    {{ if .Scaffold.user_relation -}}
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,{{ end }}
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS {{ .Scaffold.sql_table }};

-- +goose StatementEnd

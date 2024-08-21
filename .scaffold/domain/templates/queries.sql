-- name: {{ .Computed.domain_var }}ByID :one
SELECT
    *
FROM
    {{ .Scaffold.sql_table }}
WHERE
    id = $1;

-- name: {{ .Computed.domain_var }}GetAll :many
SELECT
    *
FROM
    {{ .Scaffold.sql_table }}
ORDER BY
    created_at
LIMIT
    $1 OFFSET $2;

-- name: {{ .Computed.domain_var }}GetAllCount :one
SELECT
    COUNT(*)
FROM
    {{ .Scaffold.sql_table }};

-- name: {{ .Computed.domain_var }}DeleteByID :exec
DELETE FROM
    {{ .Scaffold.sql_table }}
WHERE
    id = $1;
{{ if .Scaffold.user_relation }}
-- name: {{ .Computed.domain_var }}GetAllByUserID :many
SELECT
    *
FROM
    {{ .Scaffold.sql_table }}
WHERE
    user_id = $1
ORDER BY
    created_at
LIMIT
    $2 OFFSET $3;{{ end }}

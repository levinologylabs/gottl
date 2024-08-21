-- name: {{ .Computed.domain_var }}ByID :one
SELECT
    *
FROM
    {{ .Computed.sql_table }}
WHERE
    id = $1;


-- name: {{ .Computed.domain_var }}GetAll :many
SELECT
    *
FROM
    {{ .Computed.sql_table }}
ORDER BY
    created_at
LIMIT
    $1 OFFSET $2;

-- name: {{ .Computed.domain_var }}GetAllCount :one
SELECT
    COUNT(*)
FROM
    {{ .Computed.sql_table }};

-- name: {{ .Computed.domain_var }}DeleteByID :exec
DELETE FROM
    {{ .Computed.sql_table }}
WHERE
    id = $1;

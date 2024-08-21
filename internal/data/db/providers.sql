-- name: ProviderGetOne :one
SELECT
    id,
    created_at,
    user_id,
    provider_name,
    provider_user_id,
    metadata
FROM
    user_identity_providers
WHERE
    user_id = $1
    AND provider_name = $2
LIMIT
    1;

-- name: CreateProvider :one
INSERT INTO
    user_identity_providers (
        user_id,
        provider_name,
        provider_user_id,
        metadata
    )
VALUES
    (
        $1,
        $2,
        $3,
        COALESCE(sqlc.narg('metadata'), '{}' :: jsonb)
    ) RETURNING *;

-- name: DeleteProvider :exec
DELETE FROM
    user_identity_providers
WHERE
    user_id = $1
    AND provider_name = $2;

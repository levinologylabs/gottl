-- name: ProviderStateCreate :exec
INSERT INTO
    user_identity_provider_state (token, expires_at)
VALUES
    ($1, $2);

-- name: ProviderStateGet :one
SELECT
    *
FROM
    user_identity_provider_state
WHERE
    token = $1;

-- name: ProviderStateDelete :exec
DELETE FROM
    user_identity_provider_state
WHERE
    token = $1;

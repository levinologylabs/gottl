-- name: UserActionTokenCreate :one
INSERT INTO
    user_action_tokens (user_id, token, action, expires_at)
VALUES
    ($1, $2, $3, $4) RETURNING *;

-- name: UserActionTokenGet :one
SELECT
    *
FROM
    user_action_tokens
WHERE
    token = $1
    AND action = $2
    AND expires_at > sqlc.arg('now')
LIMIT
    1;

-- name: UserActionTokenDelete :exec
DELETE FROM
    user_action_tokens
WHERE
    id = $1;

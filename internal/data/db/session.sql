-- name: SessionCreate :exec
INSERT INTO
    user_sessions (user_id, token, expires_at)
VALUES
    ($1, $2, $3);

-- name: UserBySession :one
SELECT
    users.*
FROM
    user_sessions
    JOIN users ON user_sessions.user_id = users.id
WHERE
    user_sessions.token = $1
    AND user_sessions.expires_at > CURRENT_TIMESTAMP;

-- name: SessionDeleteByToken :exec
DELETE FROM
    user_sessions
WHERE
    token = $1;

-- name: SessionDeleteExpiredBefore :execrows
-- SessionDeleteExpiredBefore deletes every session that has expired
-- before the given time.
DELETE FROM
    user_sessions
WHERE
    expires_at < $1;

-- name: UserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: UserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: UserCreateAdmin :one
INSERT INTO
    users (username, email, password_hash, is_admin)
VALUES
    ($1, $2, $3, TRUE) RETURNING *;

-- name: UserCreate :one
INSERT INTO
    users (username, email, password_hash)
VALUES
    ($1, $2, $3) RETURNING *;

-- name: UserUpdate :one
UPDATE
    users
SET
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash)
WHERE
    id = $1 RETURNING *;

-- name: UserUpdateBilling :one
UPDATE
    users
SET
    stripe_customer_id = COALESCE(
        sqlc.narg(stripe_customer_id),
        stripe_customer_id
    ),
    stripe_subscription_id = COALESCE(
        sqlc.narg(stripe_subscription_id),
        stripe_subscription_id
    ),
    subscription_start_date = COALESCE(
        sqlc.narg(subscription_start_date),
        subscription_start_date
    ),
    subscription_ended_date = COALESCE(
        sqlc.narg(subscription_ended_date),
        subscription_ended_date
    )
WHERE
    id = $1 RETURNING *;

-- name: UserDeleteByID :exec
DELETE FROM
    users
WHERE
    id = $1;

-- name: UserGetAllCount :one
SELECT
    COUNT(*)
FROM
    users;

-- name: UserGetAll :many
SELECT
    *
FROM
    users
ORDER BY
    id
LIMIT
    $1 OFFSET $2;

-- name: UserByProvider :one
SELECT
    users.*
FROM
    users
    JOIN user_identity_providers ON users.id = user_identity_providers.user_id
WHERE
    user_identity_providers.provider_name = $1
    AND user_identity_providers.provider_user_id = $2
LIMIT
    1;

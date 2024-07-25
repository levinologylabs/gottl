-- name: CreateUser :one
INSERT INTO
    users (username, email, password_hash)
VALUES
    ($1, $2, $3) RETURNING *;

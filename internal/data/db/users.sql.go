// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const userByEmail = `-- name: UserByEmail :one
SELECT
    id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
FROM
    users
WHERE
    email = $1
`

func (q *Queries) UserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, userByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

const userByID = `-- name: UserByID :one
SELECT
    id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
FROM
    users
WHERE
    id = $1
`

func (q *Queries) UserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, userByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

const userByProvider = `-- name: UserByProvider :one
SELECT
    users.id, users.created_at, users.updated_at, users.username, users.email, users.password_hash, users.is_admin, users.stripe_customer_id, users.stripe_subscription_id, users.subscription_start_date, users.subscription_ended_date
FROM
    users
    JOIN user_identity_providers ON users.id = user_identity_providers.user_id
WHERE
    user_identity_providers.provider_name = $1
    AND user_identity_providers.provider_user_id = $2
LIMIT
    1
`

type UserByProviderParams struct {
	ProviderName   string
	ProviderUserID string
}

func (q *Queries) UserByProvider(ctx context.Context, arg UserByProviderParams) (User, error) {
	row := q.db.QueryRow(ctx, userByProvider, arg.ProviderName, arg.ProviderUserID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

const userCreate = `-- name: UserCreate :one
INSERT INTO
    users (username, email, password_hash)
VALUES
    ($1, $2, $3) RETURNING id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
`

type UserCreateParams struct {
	Username     string
	Email        string
	PasswordHash string
}

func (q *Queries) UserCreate(ctx context.Context, arg UserCreateParams) (User, error) {
	row := q.db.QueryRow(ctx, userCreate, arg.Username, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

const userCreateAdmin = `-- name: UserCreateAdmin :one
INSERT INTO
    users (username, email, password_hash, is_admin)
VALUES
    ($1, $2, $3, TRUE) RETURNING id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
`

type UserCreateAdminParams struct {
	Username     string
	Email        string
	PasswordHash string
}

func (q *Queries) UserCreateAdmin(ctx context.Context, arg UserCreateAdminParams) (User, error) {
	row := q.db.QueryRow(ctx, userCreateAdmin, arg.Username, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

const userDeleteByID = `-- name: UserDeleteByID :exec
DELETE FROM
    users
WHERE
    id = $1
`

func (q *Queries) UserDeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, userDeleteByID, id)
	return err
}

const userGetAll = `-- name: UserGetAll :many
SELECT
    id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
FROM
    users
ORDER BY
    id
LIMIT
    $1 OFFSET $2
`

type UserGetAllParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) UserGetAll(ctx context.Context, arg UserGetAllParams) ([]User, error) {
	rows, err := q.db.Query(ctx, userGetAll, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.IsAdmin,
			&i.StripeCustomerID,
			&i.StripeSubscriptionID,
			&i.SubscriptionStartDate,
			&i.SubscriptionEndedDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const userGetAllCount = `-- name: UserGetAllCount :one
SELECT
    COUNT(*)
FROM
    users
`

func (q *Queries) UserGetAllCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, userGetAllCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const userUpdate = `-- name: UserUpdate :one
UPDATE
    users
SET
    username = COALESCE($2, username),
    email = COALESCE($3, email),
    password_hash = COALESCE($4, password_hash)
WHERE
    id = $1 RETURNING id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
`

type UserUpdateParams struct {
	ID           uuid.UUID
	Username     *string
	Email        *string
	PasswordHash *string
}

func (q *Queries) UserUpdate(ctx context.Context, arg UserUpdateParams) (User, error) {
	row := q.db.QueryRow(ctx, userUpdate,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.PasswordHash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

const userUpdateBilling = `-- name: UserUpdateBilling :one
UPDATE
    users
SET
    stripe_customer_id = COALESCE(
        $2,
        stripe_customer_id
    ),
    stripe_subscription_id = COALESCE(
        $3,
        stripe_subscription_id
    ),
    subscription_start_date = COALESCE(
        $4,
        subscription_start_date
    ),
    subscription_ended_date = COALESCE(
        $5,
        subscription_ended_date
    )
WHERE
    id = $1 RETURNING id, created_at, updated_at, username, email, password_hash, is_admin, stripe_customer_id, stripe_subscription_id, subscription_start_date, subscription_ended_date
`

type UserUpdateBillingParams struct {
	ID                    uuid.UUID
	StripeCustomerID      *string
	StripeSubscriptionID  *string
	SubscriptionStartDate pgtype.Timestamp
	SubscriptionEndedDate pgtype.Timestamp
}

func (q *Queries) UserUpdateBilling(ctx context.Context, arg UserUpdateBillingParams) (User, error) {
	row := q.db.QueryRow(ctx, userUpdateBilling,
		arg.ID,
		arg.StripeCustomerID,
		arg.StripeSubscriptionID,
		arg.SubscriptionStartDate,
		arg.SubscriptionEndedDate,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.IsAdmin,
		&i.StripeCustomerID,
		&i.StripeSubscriptionID,
		&i.SubscriptionStartDate,
		&i.SubscriptionEndedDate,
	)
	return i, err
}

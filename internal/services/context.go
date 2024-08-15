package services

import (
	"context"
	"errors"

	"github.com/jalevin/gottl/internal/data/dtos"
)

var ErrNotAdmin = errors.New("user is not an admin")

type serviceCtxKey string

var (
	serviceCtxKeyUser          = serviceCtxKey("user")
	serviceCtxKeyVerifiedAdmin = serviceCtxKey("verified_admin")
)

// WithUser adds a user to the context.
func WithUser(ctx context.Context, user dtos.User) context.Context {
	return context.WithValue(ctx, serviceCtxKeyUser, user)
}

// UserFrom retrieves a user from the context.
//
// WARNING: this function will panic if the user is not found in the context.
func UserFrom(ctx context.Context) dtos.User {
	user, ok := ctx.Value(serviceCtxKeyUser).(dtos.User)
	if !ok {
		panic("user not found in context, never call UserFrom outside of an authenticated context")
	}

	return user
}

// WithVerifiedAdmin adds a verified admin to the context.
func WithVerifiedAdmin(ctx context.Context) context.Context {
	return context.WithValue(ctx, serviceCtxKeyVerifiedAdmin, true)
}

// VerifiedAdminFrom checks if the user has been identified as an admin within
// the context of the request. This is done by the middleware. If the user is
// NOT an admin, an error of ErrNotAdmin is returned.
func VerifiedAdminFrom(ctx context.Context) error {
	verified, _ := ctx.Value(serviceCtxKeyVerifiedAdmin).(bool)
	if !verified {
		return ErrNotAdmin
	}

	return nil
}

package services

import (
	"context"

	"github.com/jalevin/gottl/internal/data/dtos"
)

type serviceCtxKey string

var serviceCtxKeyUser = serviceCtxKey("user")

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

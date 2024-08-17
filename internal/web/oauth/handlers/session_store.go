package handlers

import (
	"context"

	"github.com/jalevin/gottl/internal/data/dtos"
)

// OauthSessionStore is an interface retrieving session data from oauth providers.
type OauthSessionStore interface {
	ProviderSession(ctx context.Context, providerName string, extID string, extEmail string, extName string) (dtos.UserSession, error)
}

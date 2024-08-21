package handlers

import (
	"context"

	"github.com/jalevin/gottl/internal/data/dtos"
)

// OauthSessionStore is an interface retrieving session data from oauth providers.
type OauthSessionStore interface {
	// ProviderStateUse consumes the token value if present within the store. If the token isn't found
	// in the store, this method returns an error
	ProviderStateUse(ctx context.Context, token string) error
	// ProviderStateGet retrieves a token value for managing state across the Oauth Lifecycle to ensure
	// that callbacks to the endpoints originated from this server.
	ProviderStateGet(ctx context.Context) (token string, err error)
	ProviderSession(ctx context.Context, providerName, extID, extEmail, extName string) (dtos.UserSession, error)
}

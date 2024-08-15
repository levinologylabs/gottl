package mid

import (
	"net/http"
	"strings"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/services"
)

func Authenticate(userservice *services.UserService) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			bearer = strings.TrimPrefix(bearer, "Brearer ")

			if bearer == "" {
				_ = server.Error().
					Status(http.StatusUnauthorized).
					Msg("Authorization header is required").
					Write(r.Context(), w)
				return
			}

			user, err := userservice.SessionVerify(r.Context(), bearer)
			if err != nil {
				_ = server.Error().
					Status(http.StatusUnauthorized).
					Msg("unauthorized").
					Write(r.Context(), w)
			}

			r = r.WithContext(services.WithUser(r.Context(), user))
			h.ServeHTTP(w, r)
		})
	}
}

// AuthorizeAdmin is a middleware that checks if the user within the context is
// and admin. If the user is not an admin, a 403 forbidden response is written.
// This middleware MUST come after the Authenticate middleware.
func AuthorizeAdmin() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := services.UserFrom(r.Context())

			if !user.IsAdmin {
				_ = server.Error().
					Status(http.StatusForbidden).
					Msg("forbidden").
					Write(r.Context(), w)
			}

			h.ServeHTTP(w, r)
		})
	}
}

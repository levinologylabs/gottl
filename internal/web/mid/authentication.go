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
					Msg("missing bearer token").
					Write(r.Context(), w)
				return
			}

			// TODO: Implement the token validation

			h.ServeHTTP(w, r)
		})
	}
}

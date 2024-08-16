package mid

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/services"
	"github.com/rs/zerolog"
)

type ErrorAdapter func(func(http.ResponseWriter, *http.Request) error) http.HandlerFunc

func (fn ErrorAdapter) Adapt(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return fn(next)
}

// ErrorHandler is a middleware that handles errors from downstream handlers. This 'middleware' is
// more of an adapter that wraps the next handler in a function that will handle the error returned.
// The error is logged and a response is written to the client.
//
// This does not recover from panics.
func ErrorHandler(log zerolog.Logger) ErrorAdapter {
	return func(next func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := next(w, r)

			// Early return on **SUCCESS**
			if err == nil {
				return
			}

			// Early return if the error builder was executed outside of the middleware.
			// This allows down-stream handlers or middleware to return a more specific
			// error response to the user.
			var builderResultErr server.ResponseError
			if errors.As(err, &builderResultErr) {
				cause := builderResultErr.Unwrap()

				log.Err(cause).
					Ctx(r.Context()).
					Msg("error builder executed outside of middleware")

				// No write action needed, response is already written
				return
			}

			// Use the outer error with a default message to ensure
			// that we don't leak any information to the client that
			// is embedded in the error message.
			bldr := server.Err(err).Msg("unknown error")

			switch {
			case errors.Is(err, pgx.ErrNoRows):
				bldr.Status(http.StatusNotFound).
					Msg("resource not found")
			case errors.Is(err, services.ErrNotAdmin):
				bldr.Status(http.StatusForbidden).
					Msg("forbidden")
			default:
				bldr.Status(http.StatusInternalServerError).
					Msg("internal server error")
				log.Err(err).Type("type", err).Msg("unhandled error resulted in 500 response")
			}

			err = bldr.Write(r.Context(), w)
			if err != nil {
				log.Err(err).Msg("error writing response")
			}
		}
	}
}

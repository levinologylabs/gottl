package mid

import (
	"net/http"

	"github.com/rs/zerolog"
)

type spy struct {
	http.ResponseWriter
	status int
}

func (s *spy) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func Logger(l zerolog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info().Ctx(r.Context()).Str("method", r.Method).Str("path", r.URL.Path).Msg("->")
			s := &spy{ResponseWriter: w}
			h.ServeHTTP(s, r)
			l.Info().Ctx(r.Context()).Str("method", r.Method).Str("path", r.URL.Path).Int("status", s.status).Msg("<-")
		})
	}
}

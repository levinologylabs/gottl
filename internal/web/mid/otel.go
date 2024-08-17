package mid

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

func Tracing(name string, mux *chi.Mux) func(next http.Handler) http.Handler {
	return otelchi.Middleware("gottl",
		otelchi.WithChiRoutes(mux),
		otelchi.WithTraceIDResponseHeader(nil),
	)
}

package mid

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jalevin/gottl/internal/observability/otel"
	"github.com/riandyrn/otelchi"
)

func Tracing(name string, mux *chi.Mux, os *otel.OtelService) func(next http.Handler) http.Handler {
	return otelchi.Middleware("gottl",
		otelchi.WithTracerProvider(os.TraceProvider),
		otelchi.WithChiRoutes(mux),
		otelchi.WithTraceIDResponseHeader(nil),
	)
}

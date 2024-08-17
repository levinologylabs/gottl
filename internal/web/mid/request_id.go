package mid

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jalevin/gottl/internal/core/server"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("gottl")

// TraceIDTraceHook is a zerolog hook that adds the trace ID to the log
// output.
type TraceIDTraceHook struct{}

func (h TraceIDTraceHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	_, traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return
	}

	e.Str("trace_id", traceID)
}

type reqCtxKey string

const reqCtxKeyType reqCtxKey = "trace_id"

func GetTraceIDFromContext(ctx context.Context) (context.Context, string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return ctx, traceID.String()
	}

	ctx, span := tracer.Start(ctx, "Root")
	spanCtx = trace.SpanContextFromContext(ctx)
	defer span.End()

	return ctx, spanCtx.TraceID().String()
}

func TraceID() func(http.Handler) http.Handler {
	server.SetTraceIDFunc(GetTraceIDFromContext)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tid := r.Header.Get("X-Trace-ID")
			if tid == "" {
				tid = strings.ReplaceAll(uuid.New().String(), "-", "")
			}

			w.Header().Add("X-Trace-ID", tid)

			r = r.WithContext(context.WithValue(r.Context(), reqCtxKeyType, tid))
			h.ServeHTTP(w, r)
		})
	}
}

package mid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jalevin/gottl/internal/core/server"
	"github.com/rs/zerolog"
)

// RequestIDTraceHook is a zerolog hook that adds the request ID to the log
// output.
type RequestIDTraceHook struct{}

func (h RequestIDTraceHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	requestID := getRequestIDFromContext(ctx)
	if requestID == "" {
		return
	}

	e.Str("rid", requestID)
}

type reqCtxKey string

const reqCtxKeyType reqCtxKey = "request_id"

func getRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(reqCtxKeyType).(string); ok {
		return requestID
	}
	return ""
}

func RequestID() func(http.Handler) http.Handler {
	server.SetRequestIDFunc(getRequestIDFromContext)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := r.Header.Get("X-Request-ID")
			if rid == "" {
				rid = uuid.New().String()
			}

			r = r.WithContext(context.WithValue(r.Context(), reqCtxKeyType, rid))
			h.ServeHTTP(w, r)
		})
	}
}

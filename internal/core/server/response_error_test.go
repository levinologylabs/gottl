package server

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type unwrapable interface {
	Unwrap() error
}

func unwrap(err error) error {
	asUnwrapable, ok := err.(unwrapable) //nolint:errorlint
	if !ok {
		return err
	}

	return asUnwrapable.Unwrap()
}

// UnsetTraceIDFunc sets the trace ID function to return an empty string.
// This is useful for testing.
func unsetTraceIDFunc() {
	traceIDFunc = func(ctx context.Context) (context.Context, string) {
		return context.Background(), ""
	}
}

func Test_ErrorBuilder(t *testing.T) {
	type tcase struct {
		name       string
		builder    *ErrorBuilder
		wantErr    error
		expectJSON string
		hook       func(context.Context) context.Context
	}

	cases := []tcase{
		{
			name:       "default error",
			builder:    Error(),
			wantErr:    errors.New("unknown error"),
			expectJSON: `{"message":"unknown error","statusCode":500}`,
		},
		{
			name: "message overrides error",
			builder: Error().
				Err(errors.New("test error")).
				Msg("test message"),
			wantErr:    errors.New("unknown error"),
			expectJSON: `{"message":"test message","statusCode":500}`,
		},
		{
			name: "data is included",
			builder: Error().
				Err(errors.New("test error")).
				Msg("test message").
				Data(map[string]string{"foo": "bar"}),
			wantErr:    errors.New("unknown error"),
			expectJSON: `{"message":"test message","statusCode":500,"data":{"foo":"bar"}}`,
		},
		{
			name:       "with trace ID",
			builder:    Error(),
			wantErr:    errors.New("unknown error"),
			expectJSON: `{"message":"unknown error","statusCode":500,"traceId":"test-trace-id"}`,
			hook: func(ctx context.Context) context.Context {
				tracer := otel.Tracer("gottl-test")
				fn := func(ctx context.Context) (context.Context, string) {
					ctx, span := tracer.Start(ctx, "test-span")
					defer span.End()
					spanctx := trace.SpanContextFromContext(ctx)
					return ctx, spanctx.TraceID().String()
				}

				SetTraceIDFunc(fn)

				return ctx
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer unsetTraceIDFunc()
			bg := context.Background()
			if c.hook != nil {
				bg = c.hook(bg)
			}

			writer := httptest.NewRecorder()

			err := c.builder.Write(bg, writer)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			// Unwrap
			err = unwrap(err)
			if errors.Is(err, c.wantErr) {
				t.Errorf("expected error %v, got %v", c.wantErr, err)
			}

			if writer.Code != 500 {
				t.Errorf("expected code 500, got %d", writer.Code)
			}

			if writer.Body.String() != c.expectJSON {
				t.Errorf("expected body %s, got %s", c.expectJSON, writer.Body.String())
			}
		})
	}
}

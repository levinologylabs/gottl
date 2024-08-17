package server

import "context"

type TraceIDFunc func(ctx context.Context) (context.Context, string)

var traceIDFunc TraceIDFunc = func(ctx context.Context) (context.Context, string) {
	return ctx, ""
}

// SetTraceIDFunc sets the function used to get the request ID from a context.
func SetTraceIDFunc(fn TraceIDFunc) {
	traceIDFunc = fn
}

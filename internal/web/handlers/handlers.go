// Package handlers contains the HTTP handlers for the application endpoints.
package handlers

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("handlers")

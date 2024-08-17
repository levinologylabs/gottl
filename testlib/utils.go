package testlib

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func IntegrationGuard(t *testing.T) {
	t.Helper()
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skipping integration test")
	}
}

func EnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

type tLogWriter struct {
	t *testing.T
}

func (w tLogWriter) Write(p []byte) (n int, err error) {
	// trim duplicate newline
	removed := 0
	if len(p) > 0 && p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
		removed = 1
	}

	w.t.Log(string(p))
	return len(p) + removed, nil
}

// Logger returns a new logger that is configured to write to the test's log.
func Logger(t *testing.T) zerolog.Logger {
	t.Helper()

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out: tLogWriter{t: t},
	}).With().
		Timestamp().
		Str("test", t.Name()).
		Logger()

	return logger
}

// Ptr returns a point to the provided values. Useful for initializing strings
// and other primtive types directly as pointers.
func Ptr[T any](v T) *T {
	return &v
}

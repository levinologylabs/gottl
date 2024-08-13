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
	if len(p) > 0 && p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}

	w.t.Log(string(p))
	return len(p), nil
}

// TestLogger returns a new logger child that is configured to write to the test's log.
// This is useful for writing logs to the associated test's output.
func TestLogger(t *testing.T) zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out: tLogWriter{t: t},
	}).With().
		Str("test", t.Name()).
		Logger()

	return logger
}

// Package testlib provides utilities for testing.
package testlib

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/rs/zerolog"
)

type OptionsFunc func(*options)

func randomDBName() string {
	var id [8]byte
	_, err := rand.Read(id[:])
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("test_%x", id)
}

type options struct {
	database string
}

// NewPersistent creates a new database that is persistent across tests. The
// database is created with a random name and returns a function to close the
// connection.
//
// NOTE: Databases are not cleaned up or removed. It is assumed that you are
// using a postgres instance that will be reset or cleaned up after the tests
// are run (docker, k8s, etc).
func NewPersistent(t *testing.T, tlog zerolog.Logger, fns ...OptionsFunc) *db.QueriesExt {
	t.Helper()
	IntegrationGuard(t)

	options := &options{
		database: randomDBName(),
	}
	for _, fn := range fns {
		fn(options)
	}

	dbconf := db.Config{
		Host:      EnvOrDefault("GOTTL_POSTGRES_HOST", "localhost"),
		Port:      EnvOrDefault("GOTTL_POSTGRES_PORT", "5432"),
		Username:  EnvOrDefault("GOTTL_POSTGRES_USER", "postgres"),
		Password:  EnvOrDefault("GOTTL_POSTGRES_PASSWORD", "postgres"),
		Database:  "",
		EnableSSL: false,
	}

	ctx := context.Background()

	var conn *pgx.Conn

	retries := 0
	for {
		var err error
		conn, err = pgx.Connect(ctx, dbconf.DSN())
		if err == nil {
			break
		}

		retries++
		if retries > 10 {
			t.Fatal(err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	var exists bool
	err := conn.QueryRow(ctx, "SELECT 1 FROM pg_database WHERE datname = $1", options.database).Scan(&exists)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("failed to check if database exists: %v", err)
		}

		exists = false
	}

	if !exists {
		_, err = conn.Exec(ctx, "CREATE DATABASE "+options.database)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = conn.Close(ctx)
	if err != nil {
		t.Fatal(err)
	}

	dbconf.Database = options.database

	q, err := db.NewExt(ctx, tlog, dbconf, true)
	if err != nil {
		t.Fatal(err)
	}

	return q
}

// NewDatabase creates a new database based on the shared postgres instance and
// calls the close function when the test is done using t.Cleanup. To use a shared
// instance of the database, use NewPersistent.
func NewDatabase(t *testing.T, tlog zerolog.Logger, fns ...OptionsFunc) *db.QueriesExt {
	t.Helper()

	client := NewPersistent(t, tlog, fns...)
	t.Cleanup(func() {
		_ = client.Close(context.Background())
	})

	return client
}

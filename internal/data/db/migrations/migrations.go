// Package migrations handles the database migrations using goose and embedded sql files.
package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

// AdvisoryLock is the advisory lock ID used to prevent multiple migrations
// from running concurrently. This integer was randomly selected and has no
// inherent meaning.
var AdvisoryLock = 19234

//go:embed sql
var migrationFS embed.FS

func Migrate(log zerolog.Logger, db *sql.DB) error {
	goose.SetLogger(logger{log: log})

	_ = goose.SetDialect("pgx")
	goose.SetBaseFS(migrationFS)

	if err := goose.Up(db, "sql"); err != nil {
		return fmt.Errorf("failed to run sqlite migrations: %w", err)
	}

	log.Info().Msg("successfully ran sqlite migrations")
	return nil
}

func Rollback(log zerolog.Logger, db *sql.DB) error {
	goose.SetLogger(logger{log: log})

	_ = goose.SetDialect("pgx")
	goose.SetBaseFS(migrationFS)

	if err := goose.Down(db, "sql"); err != nil {
		return fmt.Errorf("failed to rollback sqlite migrations: %w", err)
	}

	log.Info().Msg("successfully rolled back sqlite migrations")
	return nil
}

var _ goose.Logger = &logger{}

type logger struct {
	log zerolog.Logger
}

// Fatalf implements goose.Logger.
func (l logger) Fatalf(format string, v ...interface{}) {
	trimmed := strings.TrimSuffix(format, "\n")
	l.log.Fatal().Msgf(trimmed, v...)
}

// Printf implements goose.Logger.
func (l logger) Printf(format string, v ...interface{}) {
	trimmed := strings.TrimSuffix(format, "\n")
	l.log.Info().Msgf(trimmed, v...)
}

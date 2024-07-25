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

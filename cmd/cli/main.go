package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/db/migrations"
	"github.com/jalevin/gottl/internal/observability/logtools"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Postgres db.Config
	Logs     logtools.Config
}

// ConfigFromCLI parses the CLI/Config file and returns a Config struct. If the file argument is an empty string, the
// file is not read. If the file is not empty, the file is read and the Config struct is returned.
func ConfigFromCLI() (*Config, error) {
	cfg := Config{}

	const prefix = "CLI"

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			os.Exit(0)
		}
		return &cfg, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

func main() {
	cfg, err := ConfigFromCLI()
	if err != nil {
		panic(err)
	}

	log.Logger, err = logtools.New(cfg.Logs)
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		log.Fatal().Msg("missing command")
	}

	cmd := os.Args[1]

	switch cmd {
	case "migrate":
		subcmd := ""
		if len(os.Args) > 2 {
			subcmd = os.Args[2]
		}

		stdlibConn, err := sql.Open("pgx", cfg.Postgres.DSN())
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open database connection")
		}
		defer stdlibConn.Close()

		obtainlock := func() {
			_, err := stdlibConn.Exec("SELECT pg_advisory_lock($1)", migrations.AdvisoryLock)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to obtain lock")
			}

			log.Info().Msg("obtained lock")
		}

		switch subcmd {
		case "up":
			log.Info().Msg("running migrations up")
			obtainlock()

			err := migrations.Migrate(log.Logger, stdlibConn)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to run migrations up")
			}

			log.Info().Msg("successfully ran migrations up")
		case "down":
			log.Info().Msg("rolling back migrations")

			obtainlock()
			err := migrations.Rollback(log.Logger, stdlibConn)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to rollback migrations")
			}

			log.Info().Msg("successfully rolled back migrations")
		default:
			log.Fatal().Str("subcmd", subcmd).Msg("unknown subcommand")
		}

	case "seed":
		log.Info().Msg("seeding database")
	default:
		log.Fatal().Str("cmd", cmd).Msg("unknown command")
	}
}

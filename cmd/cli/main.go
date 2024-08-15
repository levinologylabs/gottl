package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/db/migrations"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/observability/logtools"
	"github.com/jalevin/gottl/internal/services"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Postgres db.Config
	Logs     logtools.Config
	Seed     struct {
		Email    string `conf:"default:admin@example.com"`
		Username string `conf:"default:admin"`
		Password string `conf:"default:admin1"`
	}
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
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {
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
		defer func() {
			conerr := stdlibConn.Close()
			log.Err(conerr).Msg("failed to close database connection")
		}()

		obtainlock := func() {
			_, err := stdlibConn.Exec("SELECT pg_advisory_lock($1)", migrations.AdvisoryLock)
			if err != nil {
				log.Error().Err(err).Msg("failed to obtain lock")
				return
			}

			log.Info().Msg("obtained lock")
		}

		switch subcmd {
		case "up":
			log.Info().Msg("running migrations up")
			obtainlock()

			err := migrations.Migrate(log.Logger, stdlibConn)
			if err != nil {
				log.Error().Err(err).Msg("failed to run migrations up")
			}

			log.Info().Msg("successfully ran migrations up")
		case "down":
			log.Info().Msg("rolling back migrations")

			obtainlock()
			err := migrations.Rollback(log.Logger, stdlibConn)
			if err != nil {
				log.Error().Err(err).Msg("failed to rollback migrations")
				return err
			}

			log.Info().Msg("successfully rolled back migrations")
		default:
			log.Error().Str("subcmd", subcmd).Msg("unknown subcommand")
			return errors.New("unknown subcommand")
		}

	case "seed":
		log.Info().Msg("seeding database")

		q, err := db.NewExt(context.Background(), log.Logger, cfg.Postgres, false)
		if err != nil {
			log.Error().Err(err).Msg("failed to create db connection")
			return err
		}

		service := services.NewUserService(log.Logger, q)

		user, err := service.Register(context.Background(), dtos.UserRegister{
			Email:    cfg.Seed.Email,
			Username: cfg.Seed.Username,
			Password: cfg.Seed.Password,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to seed database")
			return err
		}

		var (
			subCustID    = "cus_ADMIN_OVERRIDE"
			subID        = "sub_ADMIN_OVERRIDE"
			subStartDate = time.Now()
			subEndDate   = time.Now().AddDate(20, 0, 0)
		)

		_, err = service.UpdateSubscription(context.Background(), user.ID, dtos.UserUpdateSubscription{
			StripeCustomerID:      &subCustID,
			StripeSubscriptionID:  &subID,
			SubscriptionStartDate: &subStartDate,
			SubscriptionEndedDate: &subEndDate,
		})
		if err != nil {
			log.Error().Err(err).Str("step", "subscription data").Msg("failed to seed database")
			return err
		}

		log.Info().
			Str("id", user.ID.String()).
			Str("email", user.Email).
			Str("username", user.Username).
			Msg("successfully seeded database")

	default:
		log.Error().Str("cmd", cmd).Msg("unknown command")
		return errors.New("unknown command")
	}

	return nil
}

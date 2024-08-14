package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jalevin/gottl/internal/data/db/migrations"
	"github.com/rs/zerolog"
)

// QueriesExt is an extension of the generated Queries struct which also depends directly on the internal
// sql connection and allows for easier transaction handling and some basic utility methods for working
// with the database.
type QueriesExt struct {
	*Queries
	conn *pgx.Conn
}

// Close closes the connection.
func (qe *QueriesExt) Close(ctx context.Context) error {
	return qe.conn.Close(ctx)
}

// WithTx runs the given function in a transaction.
func (qe *QueriesExt) WithTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := qe.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := qe.Queries.WithTx(tx)
	if err := fn(q); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func NewExt(ctx context.Context, logger zerolog.Logger, config Config, runMigrations bool) (*QueriesExt, error) {
	var conn *pgx.Conn
	var err error

	var (
		retries = 5
		wait    = 1
		dsn     = config.DSN()
	)

	for {
		conn, err = pgx.Connect(ctx, dsn)
		if err == nil {
			err = conn.Ping(ctx)
			if err == nil {
				break
			}
		}

		if retries == 0 {
			return nil, err
		}

		retries--
		logger.Warn().Err(err).Int("retries", retries).Msg("failed to ping database, retrying...")
		time.Sleep(time.Duration(wait) * time.Second)
		wait *= 2
	}

	if runMigrations {
		var (
			retries = 5
			wait    = 5
		)

		for {
			_, err = conn.Exec(ctx, "SELECT pg_advisory_lock($1)", migrations.AdvisoryLock)
			if err != nil {
				if retries == 0 {
					return nil, err
				}

				retries--
				logger.Warn().Err(err).Int("retries", retries).Msg("failed to obtain advisory lock, retrying...")
				time.Sleep(time.Duration(wait) * time.Second)
				wait *= 2
				continue
			}

			logger.Info().Msg("obtained advisory lock for migrations")
			defer func() {
				_, err = conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", migrations.AdvisoryLock)
				if err != nil {
					logger.Error().Err(err).Msg("failed to release advisory")
				}
			}()

			break
		}

		stdlibConn, err := sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := stdlibConn.Close(); err != nil {
				logger.Error().Err(err).Msg("failed to close stdlib connection")
			}
		}()

		err = migrations.Migrate(logger, stdlibConn)
		if err != nil {
			return nil, err
		}
	}

	return &QueriesExt{
		Queries: New(conn),
		conn:    conn,
	}, nil
}

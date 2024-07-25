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
	var (
		retries = 5
		wait    = 1
		dsn     = config.DSN()
	)

	var err error

	stdlibConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer stdlibConn.Close()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	for {
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
		err := migrations.Migrate(logger, stdlibConn)
		if err != nil {
			return nil, err
		}
	}

	return &QueriesExt{
		Queries: New(conn),
		conn:    conn,
	}, nil
}

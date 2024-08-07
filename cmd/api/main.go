package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/observability/logtools"
	"github.com/jalevin/gottl/internal/web"
)

// @title                      Gottl API
// @version                    1.0
// @description                This is a standard Rest API template
// @BasePath                   /api
// @securityDefinitions.apikey Bearer
// @in                         header
// @name                       Authorization
// @description                "Type 'Bearer TOKEN' to correctly set the API Key"
func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := ConfigFromCLI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger, err := logtools.New(cfg.Logs)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apisvr := web.New(cfg.Version.Build, cfg.Web, logger)

	wg := sync.WaitGroup{}
	wg.Add(1)

	queries, err := db.NewExt(ctx, logger, cfg.Postgres, true)
	if err != nil {
		return err
	}

	defer func() {
		if err := queries.Close(ctx); err != nil {
			logger.Error().Err(err).Msg("closing queries")
		}
	}()

	go func() {
		defer wg.Done()
		err := apisvr.Start(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("server error")
		}
	}()

	wg.Wait()

	return nil
}

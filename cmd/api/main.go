package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jalevin/gottl/internal/core/mailer"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/observability/logtools"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/internal/web"
	"github.com/jalevin/gottl/internal/web/mid"
	"github.com/jalevin/gottl/internal/worker"
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

	logger, err := logtools.New(cfg.Logs, mid.RequestIDTraceHook{})
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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

	var (
		sender = mailer.NewSMTPSender(cfg.SMTP)
		wkr    = worker.New(cfg.Worker, logger, queries, sender)
		svc    = services.NewService(cfg.App, logger, queries, wkr)
		apisvr = web.New(cfg.Version.Build, cfg.Web, logger, svc)
	)

	go func() {
		defer wg.Done()
		err := apisvr.Start(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error().Err(err).Msg("server error")
		}
	}()

	go wkr.Start(ctx)

	wg.Wait()

	return nil
}

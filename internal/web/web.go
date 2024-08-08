// Package web contains the web service for the API.
package web

import (
	"context"
	"net/http"

	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/web/docs"
	"github.com/jalevin/gottl/internal/web/handlers"
	"github.com/rs/zerolog"
)

type Web struct {
	build  string
	cfg    Config
	server *http.Server
	logger zerolog.Logger
}

func New(
	build string,
	conf Config,
	logger zerolog.Logger,
) *Web {
	mux := routes(logger, build)

	w := &Web{
		build:  build,
		logger: logger,
		cfg:    conf,
		server: &http.Server{
			Handler:      mux,
			Addr:         conf.Addr(),
			IdleTimeout:  conf.IdleTimeout,
			ReadTimeout:  conf.ReadTimeout,
			WriteTimeout: conf.WriteTimeout,
		},
	}

	return w
}

// Start will start the web service and in another go routine block on the context
// when the context is done, it will shut down the server.
func (web *Web) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		web.logger.Info().Msg("shutting down web service")
		_ = web.server.Shutdown(ctx)
	}()

	web.logger.Info().Msgf("starting web service on http://%s", web.cfg.Addr())
	return web.server.ListenAndServe()
}

func routes(logger zerolog.Logger, build string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/docs/swagger.json", docs.SwaggerJSON)

	mux.HandleFunc("/api/v1/info", handlers.Info(logger, dtos.StatusResponse{
		Build: build,
	}))

	return mux
}

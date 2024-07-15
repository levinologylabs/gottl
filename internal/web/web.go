// Package web contains the web service for the API.
package web

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
)

type Web struct {
	cfg    Config
	server *http.Server
	logger zerolog.Logger
}

func New(
	conf Config,
	logger zerolog.Logger,
) *Web {
	mux := routes()

	w := &Web{
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

	if web.cfg.EnableTLS() {
		web.logger.Info().Msgf("starting web service on https://%s", web.cfg.Addr())
		return web.server.ListenAndServeTLS(web.cfg.TLSCert, web.cfg.TLSKey)
	}

	web.logger.Info().Msgf("starting web service on http://%s", web.cfg.Addr())
	return web.server.ListenAndServe()
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}

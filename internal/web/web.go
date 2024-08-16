// Package web contains the web service for the API.
package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/web/docs"
	"github.com/jalevin/gottl/internal/web/handlers"
	"github.com/jalevin/gottl/internal/web/mid"
	"github.com/riandyrn/otelchi"
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
	w := &Web{
		build:  build,
		logger: logger,
		cfg:    conf,
	}

	mux := w.routes(build)

	w.server = &http.Server{
		Handler:      mux,
		Addr:         conf.Addr(),
		IdleTimeout:  conf.IdleTimeout,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
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

func (web *Web) routes(build string) http.Handler {
	mux := chi.NewRouter()
	mux.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.CleanPath,
		middleware.StripSlashes,
		mid.RequestID(),
		otelchi.Middleware("gottl", otelchi.WithChiRoutes(mux)),
		mid.Logger(web.logger),
		middleware.AllowContentType("application/json", "text/plain", "text/html"),
	)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   web.cfg.Origins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	adapter := mid.ErrorHandler(web.logger)

	if web.cfg.EnableProfiler {
		mux.Mount("/debug", middleware.Profiler())
	}

	mux.HandleFunc("GET /docs/swagger.json", adapter.Adapt(docs.SwaggerJSON))
	mux.HandleFunc("GET /api/v1/info", adapter.Adapt(handlers.Info(dtos.StatusResponse{Build: build})))

	return mux
}

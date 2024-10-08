// Package web contains the web service for the API.
package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/observability/otel"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/internal/web/docs"
	"github.com/jalevin/gottl/internal/web/handlers"
	"github.com/jalevin/gottl/internal/web/mid"
	oauthhandler "github.com/jalevin/gottl/internal/web/oauth/handlers"
	"github.com/rs/zerolog"
)

type Web struct {
	build  string
	cfg    Config
	server *http.Server
	l      zerolog.Logger
	os     *otel.OtelService
	s      *services.Service
}

func New(
	build string,
	conf Config,
	logger zerolog.Logger,
	os *otel.OtelService,
	services *services.Service,
) *Web {
	w := &Web{
		build: build,
		l:     logger,
		cfg:   conf,
		os:    os,
		s:     services,
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
		web.l.Info().Ctx(ctx).Msg("shutting down web service")
		_ = web.server.Shutdown(ctx)
	}()

	web.l.Info().Ctx(ctx).Msgf("starting web service on http://%s", web.cfg.Addr())
	return web.server.ListenAndServe()
}

func (web *Web) routes(build string) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	if web.cfg.Otel.Enabled {
		mux.Use(mid.Tracing(web.cfg.ServiceName, mux, web.os))
	}

	mux.Use(
		middleware.RealIP,
		middleware.CleanPath,
		middleware.StripSlashes,
		mid.Logger(web.l),
		middleware.AllowContentType("application/json", "text/plain", "text/html"),
	)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   web.cfg.Origins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Trace-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	adapter := mid.ErrorHandler(web.l)

	if web.cfg.EnableProfiler {
		mux.Mount("/debug", middleware.Profiler())
	}

	mux.HandleFunc("GET /docs/swagger.json", adapter.Adapt(docs.SwaggerJSON))
	mux.HandleFunc("GET /api/v1/info", adapter.Adapt(handlers.Info(dtos.StatusResponse{Build: build})))

	userctrl := handlers.NewAuthController(web.s.Users, web.s.Passwords)

	mux.HandleFunc("POST /api/v1/users", adapter.Adapt(userctrl.Register))

	if web.cfg.Auth.IsLocalEnabled() {
		mux.HandleFunc("POST /api/v1/users/login", adapter.Adapt(userctrl.Authenticate))
		mux.HandleFunc("POST /api/v1/users/reset-password-request", adapter.Adapt(userctrl.ResetPasswordRequest))
		mux.HandleFunc("POST /api/v1/users/reset-password", adapter.Adapt(userctrl.ResetPassword))
	}

	if web.cfg.Auth.IsGoogleEnabled() {
		web.l.Info().Msg("google oauth enabled")
		google := oauthhandler.NewGoogleAuthController(web.l, web.cfg.Auth.Google, web.s.Users)

		mux.HandleFunc("GET /auth/login/google", google.Authenticate)
		mux.HandleFunc("GET /auth/callback/google", adapter.Adapt(google.Callback))
	}

	mux.Group(func(r chi.Router) {
		r.Use(mid.Authenticate(web.s.Users))

		r.HandleFunc("GET /api/v1/users/self", adapter.Adapt(userctrl.Self))
		r.HandleFunc("PATCH /api/v1/users/self", adapter.Adapt(userctrl.Update))

		// $scaffold_inject_routes
	})

	mux.Group(func(r chi.Router) {
		r.Use(
			mid.Authenticate(web.s.Users),
			mid.AuthorizeAdmin(),
		)

		admin := handlers.NewAdminController(web.s.Admin)

		r.HandleFunc("GET /api/v1/admin/users", adapter.Adapt(admin.GetAllUsers))
	})

	return mux
}

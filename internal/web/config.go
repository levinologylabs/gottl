package web

import (
	"strings"
	"time"

	"github.com/jalevin/gottl/internal/observability/otel"
	"github.com/jalevin/gottl/internal/web/oauth/providers"
)

type Config struct {
	ServiceName    string        `conf:"default:gottl"`
	Host           string        `conf:"default:0.0.0.0"`
	Port           string        `conf:"default:8080"`
	IdleTimeout    time.Duration `conf:"default:30s"`
	ReadTimeout    time.Duration `conf:"default:10s"`
	WriteTimeout   time.Duration `conf:"default:10s"`
	AllowedOrigins string        `conf:"default:http://*,http://*"`
	EnableProfiler bool          `conf:"default:false"`
	Otel           otel.Config
	Google         providers.GoogleConfig
	Auth           Auth
}

type Auth struct {
	Local  bool `conf:"default:true"`
	Google providers.GoogleConfig
}

func (a Auth) HasProvider() bool {
	return a.IsLocalEnabled() || a.IsGoogleEnabled()
}

func (a Auth) IsLocalEnabled() bool {
	return a.Local
}

func (a Auth) IsGoogleEnabled() bool {
	return a.Google.ClientID != "" && a.Google.ClientSecret != ""
}

func (cfg Config) Origins() []string {
	return strings.Split(cfg.AllowedOrigins, ",")
}

func (cfg Config) Addr() string {
	return cfg.Host + ":" + cfg.Port
}

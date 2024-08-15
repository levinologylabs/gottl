package web

import (
	"strings"
	"time"
)

type Config struct {
	Host           string        `conf:"default:0.0.0.0"`
	Port           string        `conf:"default:8080"`
	IdleTimeout    time.Duration `conf:"default:30s"`
	ReadTimeout    time.Duration `conf:"default:10s"`
	WriteTimeout   time.Duration `conf:"default:10s"`
	AllowedOrigins string        `conf:"default:http://*,http://*"`
	EnableProfiler bool          `conf:"default:false"`
}

func (cfg Config) Origins() []string {
	return strings.Split(cfg.AllowedOrigins, ",")
}

func (cfg Config) Addr() string {
	return cfg.Host + ":" + cfg.Port
}

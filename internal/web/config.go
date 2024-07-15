package web

import (
	"time"
)

type Config struct {
	Host           string        `conf:"default:0.0.0.0"`
	Port           string        `conf:"default:8080"`
	IdleTimeout    time.Duration `conf:"default:30s"`
	ReadTimeout    time.Duration `conf:"default:10s"`
	WriteTimeout   time.Duration `conf:"default:10s"`
	TLSCert        string
	TLSKey         string
}

func (cfg Config) Addr() string {
	return cfg.Host + ":" + cfg.Port
}

func (cfg Config) EnableTLS() bool {
	return cfg.TLSCert != "" && cfg.TLSKey != ""
}

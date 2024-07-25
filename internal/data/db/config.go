package db

import "fmt"

type Config struct {
	Host      string `conf:"default:localhost"`
	Port      string `conf:"default:5432"`
	Username  string `conf:"default:postgres"`
	Password  string `conf:"default:postgres"`
	Database  string `conf:"default:postgres"`
	EnableSSL bool   `conf:"default:false"`
}

func (d Config) DSN() string {
	sslMode := "disable"
	if d.EnableSSL {
		sslMode = "require"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.Username, d.Password,
		d.Host, d.Port,
		d.Database,
		sslMode,
	)
}

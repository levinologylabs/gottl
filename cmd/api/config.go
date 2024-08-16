package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/observability/logtools"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/internal/web"
	"github.com/jalevin/gottl/internal/worker"
)

type Config struct {
	conf.Version
	Web      web.Config
	Logs     logtools.Config
	Postgres db.Config
	Worker   worker.Config
	App      services.Config
}

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

// ConfigFromCLI parses the CLI/Config file and returns a Config struct. If the file argument is an empty string, the
// file is not read. If the file is not empty, the file is read and the Config struct is returned.
func ConfigFromCLI() (*Config, error) {
	cfg := Config{
		Version: conf.Version{
			Build: build(),
			Desc:  "Rest API",
		},
	}

	const prefix = "API"

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			os.Exit(0)
		}
		return &cfg, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

package logtools

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level   string `json:"level" conf:"default:debug"`
	Style   string `json:"style" conf:"default:console"`
	Color   bool   `json:"color" conf:"default:true"`
	LogFile string `json:"logFile"`
}

// New returns a new logger with the given level and style.
// This is used to simplify the creation of a logger in the main
// function of an application and keep logs consistent across the
// applications.
func New(cfg Config) (zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	var logWriter io.Writer
	switch cfg.Style {
	case "console":
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout}
	default:
		logWriter = os.Stdout
	}

	var logger zerolog.Logger
	if cfg.LogFile == "" {
		logger = zerolog.New(logWriter)
	} else {
		runLogFile, err := os.OpenFile("gottl.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			panic(err)
		}
		logger = zerolog.New(zerolog.MultiLevelWriter(logWriter, runLogFile))
	}

	logger = logger.With().
		Caller().    // adds the file and line number of the caller
		Timestamp(). // adds a timestamp to each log line
		Logger().
		Level(lvl)

	return logger, nil
}

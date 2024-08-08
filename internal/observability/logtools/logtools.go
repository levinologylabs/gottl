package logtools

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level string `json:"level" conf:"default:debug"`
	Style string `json:"style" conf:"default:console"`
	Color bool   `json:"color" conf:"default:true"`
}

// New returns a new logger with the given level and style.
// This is used to simplify the creation of a logger in the main
// function of an application and keep logs consistent across the
// applications.
func New(cfg Config, hooks ...zerolog.Hook) (zerolog.Logger, error) {
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

	l := zerolog.New(logWriter).
		With().
		Caller().    // adds the file and line number of the caller
		Timestamp(). // adds a timestamp to each log line
		Logger().
		Level(lvl).
		Hook(hooks...)

	return l, nil
}

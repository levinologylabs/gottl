package logtools

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
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

	l := zerolog.New(logWriter).
		With().
		Caller().    // adds the file and line number of the caller
		Timestamp(). // adds a timestamp to each log line
		Logger().
		Level(lvl)

	// setup loki
	// Configure the Loki hook
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithStaticLabels(lokirus.Labels{
			"app":         "example",
			"environment": "development",
		}).
		WithBasicAuth("admin", "secretpassword") // Optional

	lg := lokirus.NewLokiHookWithOpts(
		"http://localhost:3100",
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel)

	l.Hook(&LogrusLokiWrapper{Logrus: lg})

	return l, nil
}

// experimental. ship logs to loki
type LogrusLokiWrapper struct {
	Logrus *lokirus.LokiHook
}

func (hook *LogrusLokiWrapper) Run(e *zerolog.Event, level zerolog.Level, message string) {
	// TODO convert from zerolog to logrus entry
	entry := &logrus.Entry{
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Message: "TEST!",
	}

	if err := hook.Logrus.Fire(entry); err != nil {
		panic(err)
	}
}

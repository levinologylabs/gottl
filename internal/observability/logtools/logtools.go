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

	runLogFile, err := os.OpenFile("gottl.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	l := zerolog.New(zerolog.MultiLevelWriter(logWriter, runLogFile)).
		With().
		Caller().    // adds the file and line number of the caller
		Timestamp(). // adds a timestamp to each log line
		Logger().
		Level(lvl)

	//l = l.Hook(NewLokiHook())

	return l, nil
}

// Configure the Loki hook
func NewLokiHook() *LogrusLokiWrapper {
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithStaticLabels(lokirus.Labels{
			"app":         "gottl",
			"environment": "development",
		})

	lg := lokirus.NewLokiHookWithOpts(
		"http://localhost:3100",
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	)
	return &LogrusLokiWrapper{Logrus: lg}
}

// experimental. ship logs to loki
type LogrusLokiWrapper struct {
	Logrus *lokirus.LokiHook
}

func (hook *LogrusLokiWrapper) Run(e *zerolog.Event, level zerolog.Level, message string) {
	// convert zerolog level to proper logrus level
	// skipping panic and trace level as those don't
	// exist in logrus
	var ll logrus.Level
	switch level {
	case zerolog.InfoLevel:
		ll = logrus.InfoLevel
	case zerolog.DebugLevel:
		ll = logrus.DebugLevel
	case zerolog.WarnLevel:
		ll = logrus.WarnLevel
	case zerolog.ErrorLevel:
		ll = logrus.ErrorLevel
	case zerolog.FatalLevel:
		ll = logrus.FatalLevel
	default:
		ll = logrus.InfoLevel
	}

	// convert from zerolog to logrus entry
	entry := &logrus.Entry{
		Time:    time.Now(),
		Level:   ll,
		Message: message,
	}

	// maybe shouldn't panic if loki not available
	if err := hook.Logrus.Fire(entry); err != nil {
		panic(err)
	}
}

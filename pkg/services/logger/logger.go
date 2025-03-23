package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

type Config struct {
	Debug bool `env:"DEBUG"` // enables debug logging.
}

func New(
	config Config,
) *slog.Logger {
	if config.Debug {
		return NewDebug()
	}

	return NewDefault()
}

func NewDefault() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
}

//nolint:mnd
func NewDebug() *slog.Logger {
	return slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
		MaxErrorStackTrace: 10,
	}))
}

func AsPrintf(slogFunc func(string, ...any)) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		slogFunc(fmt.Sprintf(format, args...))
	}
}

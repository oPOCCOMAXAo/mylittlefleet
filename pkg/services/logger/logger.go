package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
)

type Config struct {
	Debug bool `env:"DEBUG"`
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
		Level:     slog.LevelDebug,
	}))
}

//nolint:mnd
func NewDebug() *slog.Logger {
	return slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
		MaxErrorStackTrace: 10,
	}))
}

func AsPrintf(slogFunc func(string, ...any)) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		slogFunc(fmt.Sprintf(format, args...))
	}
}

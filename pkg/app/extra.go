package app

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func NewCancelCause(
	shutdowner fx.Shutdowner,
	logger *slog.Logger,
) context.CancelCauseFunc {
	return func(err error) {
		exitCode := 0

		if err != nil {
			exitCode = 1

			logger.Error("CancelCause",
				slog.Any("error", err),
			)
		}

		err = shutdowner.Shutdown(
			fx.ExitCode(exitCode),
		)
		if err != nil {
			logger.Error("CancelCause",
				slog.Any("error", err),
			)
		}
	}
}

func NewFxLogger(
	logger *slog.Logger,
) fxevent.Logger {
	return &fxevent.SlogLogger{
		Logger: logger,
	}
}

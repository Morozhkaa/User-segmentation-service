// The logger package is responsible for initializing the logger.
package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

type LoggerOptions struct {
	IsProd bool
}

// New initializes the logger. Depending on the IsProd parameter, the logging level can be either LevelInfo or LevelDebug.
func New(opts LoggerOptions) *slog.Logger {
	switch {
	case opts.IsProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}

// Get returns the previously initialized logger.
func Get() *slog.Logger {
	if log == nil {
		return New(LoggerOptions{IsProd: false})
	}
	return log
}

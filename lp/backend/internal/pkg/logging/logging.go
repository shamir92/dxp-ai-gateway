// Package logging provides the process-wide structured logger.
//
// Everything goes to stdout as JSON so the platform log pipeline can index it.
// This matters because the backend previously used the stdlib logger
// (free-form text) and, more importantly, errors were discarded before anything
// could log them.
package logging

import (
	"log/slog"
	"os"
	"sync/atomic"
)

// level holds the active threshold. It is set from LOG_LEVEL by Setup, which
// runs after .env is loaded — reading the environment at package init would miss
// a level configured in .env.
var level atomic.Int32

func init() {
	level.Store(int32(slog.LevelInfo))
}

// Setup installs the JSON logger as the process default and applies LOG_LEVEL.
// Call it once, after configuration is loaded and before serving.
func Setup() {
	level.Store(int32(parseLevel(os.Getenv("LOG_LEVEL"))))
	slog.SetDefault(Logger())
}

// Logger returns the process logger. It is safe for concurrent use.
func Logger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(level.Load()),
	}))
}

func parseLevel(raw string) slog.Level {
	switch raw {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

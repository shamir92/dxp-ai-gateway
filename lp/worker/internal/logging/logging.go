// Package logging provides the worker's process-wide structured logger.
//
// Everything goes to stdout as JSON so the worker's output has the same shape as
// the backend's. The worker previously used the stdlib default handler (text to
// stderr), which meant a startup line could not be parsed by the same pipeline
// that reads the API's logs.
//
// This is a deliberate copy of the backend's logging package: the worker is a
// standalone module and must not import backend code.
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
// Call it once, after .env is loaded and before the first log line is expected,
// so that even a startup failure is emitted as JSON.
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

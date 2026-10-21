// Package observability initialises Sentry for the worker.
//
// Optional by design: with no SENTRY_DSN the SDK stays disabled, so the worker
// runs locally without network calls. The scrubbers matter more here than in the
// API because the worker handles decrypted site passwords.
package observability

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

var sensitiveHeaders = map[string]struct{}{
	"cookie":              {},
	"set-cookie":          {},
	"authorization":       {},
	"x-security-proof":    {},
	"proxy-authorization": {},
}

// Init configures Sentry from the environment and reports whether it is active.
func Init() bool {
	dsn := strings.TrimSpace(os.Getenv("SENTRY_DSN"))
	if dsn == "" {
		slog.Info("sentry disabled: SENTRY_DSN not set")
		return false
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      envOr("SENTRY_ENVIRONMENT", "development"),
		Release:          resolvedRelease(),
		TracesSampleRate: tracesSampleRate(),
		AttachStacktrace: true,
		BeforeSend:       scrubEvent,
	})
	if err != nil {
		slog.Error("sentry init failed", slog.Any("error", err))
		return false
	}
	slog.Info("sentry enabled",
		slog.String("environment", envOr("SENTRY_ENVIRONMENT", "development")),
		slog.String("release", resolvedRelease()),
	)
	return true
}

func tracesSampleRate() float64 {
	raw := strings.TrimSpace(os.Getenv("SENTRY_TRACES_SAMPLE_RATE"))
	if raw == "" {
		return 0.1
	}
	rate, err := strconv.ParseFloat(raw, 64)
	if err != nil || rate < 0 || rate > 1 {
		return 0.1
	}
	return rate
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// scrubEvent strips credentials. The worker's payloads include a site password,
// so a request body must never be serialised into an event.
func scrubEvent(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	if event.Request != nil {
		event.Request.Cookies = ""
		event.Request.Data = ""
		// Match case-insensitively: Sentry stores header names in canonical
		// form ("Cookie"), so a lowercase-only lookup would leave them in.
		for name := range event.Request.Headers {
			if _, ok := sensitiveHeaders[strings.ToLower(name)]; ok {
				delete(event.Request.Headers, name)
			}
		}
	}
	if event.User.Email != "" {
		event.User.Email = ""
	}
	return event
}

// Flush waits briefly for queued events.
func Flush() {
	sentry.Flush(2 * time.Second)
}

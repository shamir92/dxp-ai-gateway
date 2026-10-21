// Package observability initialises Sentry for error capture and tracing.
//
// Sentry is entirely optional: with no SENTRY_DSN the SDK is left disabled, so
// local development and `vercel dev` run without network calls or warnings.
package observability

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

// sensitiveHeaders are dropped from captured events. Cookies carry the session
// token and Authorization carries bearer tokens (including the "nap_…" personal
// access tokens the gateway mints for external sites).
var sensitiveHeaders = map[string]struct{}{
	"cookie":              {},
	"set-cookie":          {},
	"authorization":       {},
	"x-security-proof":    {},
	"x-csrf-protection":   {},
	"proxy-authorization": {},
}

// secretPrefixes are value shapes that must never reach Sentry even when they
// appear somewhere unexpected (a URL, a message, a breadcrumb payload).
var secretPrefixes = []string{"nap_", "sk-", "gho_", "ghp_", "Bearer "}

// Init configures Sentry from the environment and reports whether it is active.
//
// The DSN is read from SENTRY_DSN and defaults to disabled. The release is
// resolved by resolvedRelease: SENTRY_RELEASE if set, otherwise the platform's
// commit sha, so an issue can be tied to a build without extra configuration.
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
		// Attach stack traces to plain CaptureMessage calls too, otherwise a
		// worker retry notice arrives with no location at all.
		AttachStacktrace: true,
		BeforeSend:       scrubEvent,
		BeforeBreadcrumb: scrubBreadcrumb,
	})
	if err != nil {
		// A bad DSN must not take the service down.
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
		// Tracing costs quota, so sample rather than capture everything.
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

// scrubEvent removes credentials from an event before it leaves the process.
func scrubEvent(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	if event.Request != nil {
		stripSensitiveHeaders(event.Request.Headers)
		event.Request.Cookies = ""
		event.Request.Data = ""
		event.Request.QueryString = redactQuery(event.Request.QueryString)
	}
	if event.User.Email != "" {
		// Keep the account identifiable without shipping the address; the ID is
		// left intact so a user's issues still group together.
		event.User.Email = ""
	}
	return event
}

// stripSensitiveHeaders deletes credential headers in place, matching
// case-insensitively: Sentry stores header names in canonical form ("Cookie"),
// so a lowercase-only comparison silently leaves them in the event.
func stripSensitiveHeaders(headers map[string]string) {
	for name := range headers {
		if _, ok := sensitiveHeaders[strings.ToLower(name)]; ok {
			delete(headers, name)
		}
	}
}

// scrubBreadcrumb drops HTTP breadcrumbs' payloads and any value carrying a
// recognisable secret prefix.
func scrubBreadcrumb(crumb *sentry.Breadcrumb, _ *sentry.BreadcrumbHint) *sentry.Breadcrumb {
	if crumb.Data != nil {
		for name := range crumb.Data {
			if _, ok := sensitiveHeaders[strings.ToLower(name)]; ok {
				delete(crumb.Data, name)
			}
		}
		if body, ok := crumb.Data["body"]; ok && looksSecret(toText(body)) {
			delete(crumb.Data, "body")
		}
	}
	if looksSecret(crumb.Message) {
		crumb.Message = "[redacted]"
	}
	return crumb
}

func toText(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func looksSecret(s string) bool {
	if s == "" {
		return false
	}
	for _, prefix := range secretPrefixes {
		if strings.Contains(s, prefix) {
			return true
		}
	}
	return false
}

// redactQuery strips the value of any query parameter that looks like a
// credential, keeping the key so the request shape is still readable.
func redactQuery(raw string) string {
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, "&")
	for i, part := range parts {
		key, value, found := strings.Cut(part, "=")
		if !found {
			continue
		}
		if looksSecret(value) || strings.Contains(strings.ToLower(key), "token") || strings.Contains(strings.ToLower(key), "secret") {
			parts[i] = key + "=[redacted]"
		}
	}
	return strings.Join(parts, "&")
}

// Flush waits briefly for queued events, for use on shutdown.
func Flush() {
	sentry.Flush(2 * time.Second)
}

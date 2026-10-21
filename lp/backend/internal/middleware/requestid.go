package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/pkg/logging"
)

const (
	// ContextLogger holds the request-scoped *slog.Logger.
	ContextLogger = "logger"
	// ContextRequestID holds the correlation ID for the request.
	ContextRequestID = "requestID"

	// HeaderRequestID is echoed back so a user-reported failure can be traced
	// to one request without guessing from timestamps.
	HeaderRequestID = "X-Request-ID"
)

// RequestID assigns (or adopts) a correlation ID and installs a logger carrying
// it, so every line logged while handling this request can be joined up.
//
// An inbound X-Request-ID is trusted only up to a sane shape: it comes from the
// caller, and echoing arbitrary bytes into logs (and response headers) would let
// someone forge or split log lines.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := sanitizeRequestID(c.GetHeader(HeaderRequestID))
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(ContextRequestID, id)
		c.Header(HeaderRequestID, id)

		logger := logging.Logger().With(
			slog.String("request_id", id),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)
		c.Set(ContextLogger, logger)

		c.Next()
	}
}

// sanitizeRequestID keeps only an obviously-safe token, bounded in length.
func sanitizeRequestID(raw string) string {
	const maxLen = 64
	if raw == "" || len(raw) > maxLen {
		return ""
	}
	for i := 0; i < len(raw); i++ {
		switch ch := raw[i]; {
		case ch >= 'a' && ch <= 'z', ch >= 'A' && ch <= 'Z', ch >= '0' && ch <= '9', ch == '-', ch == '_', ch == '.':
		default:
			return ""
		}
	}
	return raw
}

// LoggerFrom returns the request-scoped logger, falling back to the process
// logger when the middleware did not run (e.g. an early panic).
func LoggerFrom(c *gin.Context) *slog.Logger {
	if v, ok := c.Get(ContextLogger); ok {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return logging.Logger()
}

// RequestIDFrom returns the correlation ID, if one was assigned.
func RequestIDFrom(c *gin.Context) string {
	if v, ok := c.Get(ContextRequestID); ok {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger emits one structured line per request, replacing gin.Logger().
//
// Level is chosen by status so that alerting can key off level alone: 5xx is an
// incident, 4xx is client behaviour, and successful traffic stays at INFO.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			// Record the shape of the query, but never what it carried: OAuth
			// callbacks put codes and state here.
			path += "?" + raw
		}

		c.Next()

		status := c.Writer.Status()
		log := LoggerFrom(c).With(
			slog.Int("status", status),
			slog.Duration("latency", time.Since(start)),
			slog.Int("bytes", c.Writer.Size()),
			slog.String("url", path),
			slog.String("ip", c.ClientIP()),
		)

		switch {
		case status >= 500:
			log.Error("request failed")
		case status >= 400:
			log.Warn("request rejected")
		default:
			log.Info("request completed")
		}
	}
}

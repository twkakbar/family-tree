package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger emits one structured log line per request. Health probes are logged
// at debug level to keep production logs readable.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"request_id", RequestIDFromContext(c),
		}
		if userID, ok := CurrentUserID(c); ok {
			attrs = append(attrs, "user_id", userID.String())
		}

		switch {
		case path == "/healthz" || path == "/readyz":
			slog.Debug("request", attrs...)
		case c.Writer.Status() >= 500:
			slog.Error("request", attrs...)
		case c.Writer.Status() >= 400:
			slog.Warn("request", attrs...)
		default:
			slog.Info("request", attrs...)
		}
	}
}

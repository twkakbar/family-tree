// Package middleware contains cross-cutting HTTP concerns: request identity,
// logging, panic recovery, security headers, body limits, CORS and
// authentication/authorization. Middleware never contains business logic.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "request_id"
	// maxRequestIDLen guards against header abuse when honoring client ids.
	maxRequestIDLen = 64
)

// RequestID assigns every request a unique id (honoring a sane inbound
// X-Request-ID) and reflects it in the response for log correlation.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" || len(id) > maxRequestIDLen {
			id = uuid.NewString()
		}
		c.Set(requestIDKey, id)
		c.Writer.Header().Set(requestIDHeader, id)
		c.Next()
	}
}

// RequestIDFromContext returns the request id set by RequestID ("" if unset).
func RequestIDFromContext(c *gin.Context) string {
	return c.GetString(requestIDKey)
}

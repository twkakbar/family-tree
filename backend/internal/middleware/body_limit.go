package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit caps request body size (bytes). Oversized bodies fail on read
// with http.MaxBytesReader, which Gin's binding surfaces as a 400.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

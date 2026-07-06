package middleware

import "github.com/gin-gonic/gin"

// SecureHeaders sets defense-in-depth response headers. Nginx sets them too
// in production; the API remains safe when accessed directly (dev, probes).
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

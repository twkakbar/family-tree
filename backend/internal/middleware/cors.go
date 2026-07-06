package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS implements a minimal allowlist-based CORS policy with credentials
// support. Production is same-origin behind nginx, so this mainly serves
// hybrid development (Next dev server on :3000 talking to the API on :8080
// directly). A dedicated library is unnecessary at this size.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[strings.TrimRight(o, "/")] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}
		if _, ok := allowed[strings.TrimRight(origin, "/")]; !ok {
			// Not an allowed cross-origin caller: continue without CORS
			// headers; the browser will block the response.
			c.Next()
			return
		}

		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			h.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
			h.Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

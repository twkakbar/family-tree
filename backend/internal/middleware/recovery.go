package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// Recovery converts panics into logged 500 responses using the standard
// error envelope, so clients never see a broken or HTML error page.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		slog.Error("panic recovered",
			"panic", recovered,
			"path", c.Request.URL.Path,
			"request_id", RequestIDFromContext(c),
		)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	})
}

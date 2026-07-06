// Package routes wires middleware and handlers into the Gin engine. It owns
// URL structure and middleware ordering — nothing else.
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/twkakbar/family-tree/backend/internal/config"
	"github.com/twkakbar/family-tree/backend/internal/handler"
	"github.com/twkakbar/family-tree/backend/internal/middleware"
	"github.com/twkakbar/family-tree/backend/pkg/response"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

// Handlers aggregates all HTTP handlers for route registration.
type Handlers struct {
	Auth   *handler.AuthHandler
	Users  *handler.UserHandler
	Health *handler.HealthHandler
}

// New builds the fully configured Gin engine.
func New(cfg *config.Config, jwt *token.JWTManager, h Handlers) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// Client IPs come from nginx's X-Forwarded-For; only trust configured proxies.
	_ = r.SetTrustedProxies(cfg.TrustedProxies)

	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.SecureHeaders(),
		middleware.BodyLimit(cfg.Upload.MaxUploadMB*1024*1024),
		middleware.CORS(cfg.CORS.AllowedOrigins),
	)

	// Probes (unversioned; used by container healthchecks).
	r.GET("/healthz", h.Health.Live)
	r.GET("/readyz", h.Health.Ready)

	// Uploaded files. In production nginx serves this path from the shared
	// volume; this route covers development and direct API access.
	r.Static("/uploads", cfg.Upload.Dir)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.Refresh)
			auth.POST("/logout", h.Auth.Logout)
		}

		users := v1.Group("/users", middleware.RequireAuth(jwt))
		{
			users.GET("/me", h.Users.GetMe)
			users.PATCH("/me", h.Users.UpdateMe)
		}

		// Admin-only namespace (user management endpoints land in Phase 6).
		_ = v1.Group("/admin", middleware.RequireAuth(jwt), middleware.RequireRole("admin"))
	}

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "route not found")
	})
	r.NoMethod(func(c *gin.Context) {
		response.Error(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	})

	return r
}

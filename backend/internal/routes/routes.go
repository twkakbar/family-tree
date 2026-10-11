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
	Auth           *handler.AuthHandler
	Users          *handler.UserHandler
	Families       *handler.FamilyHandler
	Persons        *handler.PersonHandler
	Public         *handler.PublicHandler
	AccessRequests *handler.AccessRequestHandler
	Transfer       *handler.TransferHandler
	AdminUsers     *handler.AdminUserHandler
	Health         *handler.HealthHandler
}

// New builds the fully configured Gin engine.
func New(cfg *config.Config, jwt *token.JWTManager, h Handlers) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	_ = r.SetTrustedProxies(cfg.TrustedProxies)

	r.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.SecureHeaders(),
		middleware.BodyLimit(cfg.Upload.MaxUploadMB*1024*1024),
		middleware.CORS(cfg.CORS.AllowedOrigins),
	)

	r.GET("/healthz", h.Health.Live)
	r.GET("/readyz", h.Health.Ready)
	r.Static("/uploads", cfg.Upload.Dir)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Auth.Register)
			auth.POST("/login", h.Auth.Login)
			auth.POST("/refresh", h.Auth.Refresh)
			auth.POST("/logout", h.Auth.Logout)
			auth.POST("/forgot-password", h.Auth.ForgotPassword)
			auth.POST("/reset-password", h.Auth.ResetPassword)
			auth.POST("/verify-email", h.Auth.VerifyEmail)
			auth.POST("/resend-verification", h.Auth.ResendVerification)
		}

		public := v1.Group("/public")
		{
			public.GET("/meta", h.Public.GetMeta)
			public.GET("/tree", h.Public.GetTree)
			public.GET("/persons/:id", h.Public.GetPerson)
			public.GET("/persons/:id/ancestors", h.Public.GetAncestors)
			public.GET("/persons/:id/descendants", h.Public.GetDescendants)
		}

		users := v1.Group("/users", middleware.RequireAuth(jwt))
		{
			users.GET("/me", h.Users.GetMe)
			users.PATCH("/me", h.Users.UpdateMe)
		}

		access := v1.Group("/access-requests", middleware.RequireAuth(jwt))
		{
			access.POST("", h.AccessRequests.Submit)
			access.GET("/me", h.AccessRequests.MyLatest)
		}

		studio := v1.Group("/studio", middleware.RequireAuth(jwt))
		{
			studio.GET("/can-edit", h.AccessRequests.CanEdit)
		}

		// Authenticated mutations still use family/person routes (studio UI
		// resolves the royal family id via /studio/can-edit).
		families := v1.Group("/families", middleware.RequireAuth(jwt))
		{
			families.GET("/:id", h.Families.GetFamily)
			families.GET("/:id/export", h.Transfer.ExportFamily)
			families.POST("/:id/import", h.Transfer.ImportFamily)
			families.GET("/:id/persons", h.Persons.ListPersons)
			families.POST("/:id/persons", h.Persons.CreatePerson)
			families.GET("/:id/tree", h.Persons.GetFamilyTree)
			families.POST("/:id/relationships", h.Persons.CreateRelationship)
		}

		persons := v1.Group("/persons", middleware.RequireAuth(jwt))
		{
			persons.GET("/:id", h.Persons.GetPerson)
			persons.PATCH("/:id", h.Persons.UpdatePerson)
			persons.PUT("/:id/photo", h.Persons.UploadPhoto)
			persons.DELETE("/:id", h.Persons.DeletePerson)
			persons.GET("/:id/ancestors", h.Persons.GetAncestors)
			persons.GET("/:id/descendants", h.Persons.GetDescendants)
		}

		relationships := v1.Group("/relationships", middleware.RequireAuth(jwt))
		{
			relationships.DELETE("/:id", h.Persons.DeleteRelationship)
		}

		admin := v1.Group("/admin", middleware.RequireAuth(jwt), middleware.RequireRole("admin"))
		{
			admin.GET("/access-requests", h.AccessRequests.ListPending)
			admin.POST("/access-requests/:id/approve", h.AccessRequests.Approve)
			admin.POST("/access-requests/:id/reject", h.AccessRequests.Reject)
			admin.GET("/users", h.AdminUsers.ListUsers)
			admin.PATCH("/users/:id", h.AdminUsers.UpdateUser)
			admin.DELETE("/users/:id", h.AdminUsers.DeleteUser)
			admin.POST("/users/:id/password-reset", h.AdminUsers.IssuePasswordReset)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "route not found")
	})
	r.NoMethod(func(c *gin.Context) {
		response.Error(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	})

	return r
}

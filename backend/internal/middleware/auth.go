package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/pkg/response"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

// Context keys set by RequireAuth.
const (
	ctxUserIDKey = "auth_user_id"
	ctxRoleKey   = "auth_role"
)

// RequireAuth validates the Bearer access token and stores the caller's
// identity in the request context. Stateless by design: role/revocation
// changes propagate within the access token TTL (15 min default).
func RequireAuth(jwt *token.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		scheme, credentials, found := strings.Cut(authz, " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || credentials == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
			return
		}

		claims, err := jwt.Verify(strings.TrimSpace(credentials))
		if err != nil {
			if errors.Is(err, token.ErrExpired) {
				response.Error(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "token has expired")
				return
			}
			response.Error(c, http.StatusUnauthorized, "TOKEN_INVALID", "token is invalid")
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxRoleKey, claims.Role)
		c.Next()
	}
}

// RequireRole allows only callers whose global role is in the allowlist.
// Must be registered after RequireAuth.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role := c.GetString(ctxRoleKey)
		if _, ok := allowed[role]; !ok {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "permission denied")
			return
		}
		c.Next()
	}
}

// CurrentUserID returns the authenticated caller's id, when present.
func CurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get(ctxUserIDKey)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

// CurrentRole returns the authenticated caller's global role ("" if absent).
func CurrentRole(c *gin.Context) string {
	return c.GetString(ctxRoleKey)
}

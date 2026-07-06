package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/twkakbar/family-tree/backend/internal/config"
	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/response"
)

// refreshCookieName is the httpOnly cookie carrying the opaque refresh token.
const refreshCookieName = "ft_refresh"

// AuthHandler adapts HTTP to the AuthService.
type AuthHandler struct {
	auth domain.AuthService
	cfg  config.AuthConfig
}

// NewAuthHandler builds the auth endpoints handler.
func NewAuthHandler(auth domain.AuthService, cfg config.AuthConfig) *AuthHandler {
	return &AuthHandler{auth: auth, cfg: cfg}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.auth.Register(c.Request.Context(), domain.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshTokenExpiresAt)
	response.OK(c, http.StatusCreated, toAuthResponse(result))
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.auth.Login(c.Request.Context(), domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshTokenExpiresAt)
	response.OK(c, http.StatusOK, toAuthResponse(result))
}

// Refresh handles POST /api/v1/auth/refresh. The refresh token is read from
// the httpOnly cookie; a new rotated cookie is set on success.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, _ := c.Cookie(refreshCookieName)

	result, err := h.auth.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		h.clearRefreshCookie(c)
		respondError(c, err)
		return
	}

	h.setRefreshCookie(c, result.RefreshToken, result.RefreshTokenExpiresAt)
	response.OK(c, http.StatusOK, toAuthResponse(result))
}

// Logout handles POST /api/v1/auth/logout. Idempotent: always clears the
// cookie and returns 204, even when no valid session existed.
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie(refreshCookieName)

	if err := h.auth.Logout(c.Request.Context(), refreshToken); err != nil {
		respondError(c, err)
		return
	}

	h.clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string, expiresAt time.Time) {
	c.SetSameSite(http.SameSiteLaxMode)
	maxAge := int(time.Until(expiresAt).Seconds())
	c.SetCookie(refreshCookieName, token, maxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(refreshCookieName, "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

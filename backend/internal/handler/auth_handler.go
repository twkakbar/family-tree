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
	auth             domain.AuthService
	cfg              config.AuthConfig
	exposeResetToken bool // non-production: include token in forgot-password response
}

// NewAuthHandler builds the auth endpoints handler.
func NewAuthHandler(auth domain.AuthService, cfg config.AuthConfig, exposeResetToken bool) *AuthHandler {
	return &AuthHandler{auth: auth, cfg: cfg, exposeResetToken: exposeResetToken}
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
	if result.VerificationRequired || result.Session == nil {
		response.OK(c, http.StatusCreated, gin.H{
			"verification_required": true,
			"message":               "Check your email to verify your account before signing in.",
		})
		return
	}

	h.setRefreshCookie(c, result.Session.RefreshToken, result.Session.RefreshTokenExpiresAt)
	response.OK(c, http.StatusCreated, toAuthResponse(result.Session))
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

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,strong_password"`
}

// ForgotPassword handles POST /api/v1/auth/forgot-password.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if !bindJSON(c, &req) {
		return
	}
	raw, err := h.auth.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		respondError(c, err)
		return
	}
	payload := gin.H{
		"requested": true,
		"message":   "If that email is registered, a reset token was created.",
	}
	if h.exposeResetToken && raw != "" {
		payload["reset_token"] = raw
	}
	response.OK(c, http.StatusOK, payload)
}

// ResetPassword handles POST /api/v1/auth/reset-password.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.auth.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"reset": true})
}

type verifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type resendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmail handles POST /api/v1/auth/verify-email.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req verifyEmailRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.auth.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{"verified": true})
}

// ResendVerification handles POST /api/v1/auth/resend-verification.
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req resendVerificationRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := h.auth.ResendVerification(c.Request.Context(), req.Email); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, http.StatusOK, gin.H{
		"sent":    true,
		"message": "If that account still needs verification, a new email was sent.",
	})
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

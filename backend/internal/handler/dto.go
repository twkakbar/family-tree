// Package handler contains the HTTP adapters: request/response DTOs,
// Gin handlers that delegate to domain services, and the central mapping
// from domain errors to HTTP responses. No business logic lives here.
package handler

import (
	"time"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

// ---- Requests ------------------------------------------------------------

// RegisterRequest creates an account.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,strong_password"`
}

// LoginRequest authenticates with credentials.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateMeRequest updates the caller's profile.
type UpdateMeRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// ---- Responses -----------------------------------------------------------

// UserResponse is the public representation of a user account.
type UserResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

func toUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:            u.ID.String(),
		Name:          u.Name,
		Email:         u.Email,
		Role:          string(u.Role),
		EmailVerified: u.EmailVerified(),
		CreatedAt:     u.CreatedAt,
	}
}

// AuthResponse is returned by register/login/refresh. The refresh token is
// NOT part of the body; it travels only in the httpOnly cookie.
type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"` // seconds until access token expiry
	User        UserResponse `json:"user"`
}

func toAuthResponse(r *domain.AuthResult) AuthResponse {
	return AuthResponse{
		AccessToken: r.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(r.AccessTokenExpiresAt).Seconds()),
		User:        toUserResponse(r.User),
	}
}

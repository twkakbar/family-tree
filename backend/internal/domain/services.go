package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RegisterInput carries validated data for creating an account.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// LoginInput carries credentials for authentication.
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult is returned by every operation that establishes a session.
// The raw RefreshToken is delivered to the client exactly once (cookie);
// only its hash is stored.
type AuthResult struct {
	User                  *User
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

// AuthService implements registration and the token lifecycle
// (login, refresh rotation with reuse detection, logout).
type AuthService interface {
	Register(ctx context.Context, in RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, in LoginInput) (*AuthResult, error)
	// Refresh rotates a valid refresh token into a new token pair. Presenting
	// a revoked token revokes every session of that user (reuse detection).
	Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	// Logout revokes the presented refresh token. It is idempotent.
	Logout(ctx context.Context, refreshToken string) error
	// EnsureAdmin creates the bootstrap admin account when credentials are
	// configured and no account with that email exists. No-op otherwise.
	EnsureAdmin(ctx context.Context, name, email, password string) error
}

// UpdateProfileInput carries the mutable fields of the caller's own profile.
type UpdateProfileInput struct {
	Name string
}

// UserService exposes account profile operations.
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, in UpdateProfileInput) (*User, error)
}

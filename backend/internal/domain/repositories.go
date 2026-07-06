package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository persists User entities.
// Implementations translate driver errors into domain errors
// (ErrNotFound, ErrEmailTaken) and must respect soft deletion.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

// RefreshTokenRepository persists refresh token hashes.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	// GetByHash returns the token row for a SHA-256 hash regardless of its
	// revocation state; callers decide how to treat revoked tokens.
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

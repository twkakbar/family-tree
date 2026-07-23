package domain

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken is an opaque one-time password reset credential.
// Only the hash is stored.
type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// TableName maps to the migration table.
func (PasswordResetToken) TableName() string { return "password_reset_tokens" }

// Valid reports whether the token can still be consumed.
func (t *PasswordResetToken) Valid(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

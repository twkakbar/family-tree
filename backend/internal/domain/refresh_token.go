package domain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken persists the SHA-256 hash of an opaque refresh token.
// The raw token is only ever held by the client (httpOnly cookie).
// Rows are hard-deleted/revoked, never soft-deleted.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

// Active reports whether the token can still be redeemed at the given time.
func (t *RefreshToken) Active(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

// Revoked reports whether the token has been explicitly revoked.
func (t *RefreshToken) Revoked() bool { return t.RevokedAt != nil }

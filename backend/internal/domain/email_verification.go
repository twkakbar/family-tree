package domain

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerificationToken is a one-time link credential. Only the hash is stored.
type EmailVerificationToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// TableName maps to the migration table.
func (EmailVerificationToken) TableName() string { return "email_verification_tokens" }

// Valid reports whether the token can still be consumed.
func (t *EmailVerificationToken) Valid(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

// EmailDispatch records one outbound email so sends can be rate-limited.
type EmailDispatch struct {
	ID      uuid.UUID
	Email   string
	Purpose string
	SentAt  time.Time
}

// TableName maps to the migration table.
func (EmailDispatch) TableName() string { return "email_dispatches" }

// Verification purpose stored in email_dispatches.purpose.
const EmailPurposeVerify = "verify"

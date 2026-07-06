// Package password wraps bcrypt hashing behind a minimal, mockable API.
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hasher hashes and verifies passwords with bcrypt.
type Hasher struct {
	cost int
}

// NewHasher builds a Hasher. Costs outside bcrypt's valid range fall back to
// bcrypt.DefaultCost (tests use low costs for speed; production uses 12).
func NewHasher(cost int) *Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &Hasher{cost: cost}
}

// Hash returns the bcrypt hash of a plaintext password.
func (h *Hasher) Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(hashed), nil
}

// Verify reports whether plain matches the stored bcrypt hash.
func (h *Hasher) Verify(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

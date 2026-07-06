package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// opaqueTokenBytes is the entropy of generated opaque tokens (256 bits).
const opaqueTokenBytes = 32

// NewOpaque returns a cryptographically random, URL-safe opaque token.
// Used for refresh tokens; only its hash may be persisted.
func NewOpaque() (string, error) {
	buf := make([]byte, opaqueTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating opaque token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashOpaque returns the hex-encoded SHA-256 digest of an opaque token.
// SHA-256 (not bcrypt) is appropriate here: the input already has 256 bits
// of entropy, so brute-forcing the digest is infeasible and lookups stay O(1).
func HashOpaque(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

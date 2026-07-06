// Package token provides JWT access-token management and opaque
// refresh-token generation. It is domain-agnostic and reusable.
package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Exported errors so callers can branch without depending on jwt internals.
var (
	// ErrExpired signals a structurally valid but expired token.
	ErrExpired = errors.New("token expired")
	// ErrInvalid signals a malformed, forged or otherwise unusable token.
	ErrInvalid = errors.New("token invalid")
)

// Claims are the JWT claims embedded in access tokens.
type Claims struct {
	UserID uuid.UUID `json:"uid"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager signs and verifies HS256 access tokens.
type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

// NewJWTManager builds a manager. The secret must be at least 32 bytes;
// config validation enforces this before construction.
func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

// TTL returns the configured access-token lifetime.
func (m *JWTManager) TTL() time.Duration { return m.ttl }

// Generate signs a new access token for the given user.
// It returns the compact token and its expiry time.
func (m *JWTManager) Generate(userID uuid.UUID, role string) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID.String(),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("signing access token: %w", err)
	}
	return signed, expiresAt, nil
}

// Verify parses and validates a compact access token.
func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, ErrInvalid
	}
	if claims.UserID == uuid.Nil {
		return nil, ErrInvalid
	}
	return claims, nil
}

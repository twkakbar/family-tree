// Package service contains ALL business logic. Services implement the
// interfaces declared in domain, depend only on domain repository interfaces
// and pkg helpers, and are unit-tested with in-memory fakes.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/password"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

// authService implements domain.AuthService.
type authService struct {
	users           domain.UserRepository
	refreshTokens   domain.RefreshTokenRepository
	jwt             *token.JWTManager
	hasher          *password.Hasher
	refreshTokenTTL time.Duration
}

// NewAuthService wires the authentication business logic.
func NewAuthService(
	users domain.UserRepository,
	refreshTokens domain.RefreshTokenRepository,
	jwt *token.JWTManager,
	hasher *password.Hasher,
	refreshTokenTTL time.Duration,
) domain.AuthService {
	return &authService{
		users:           users,
		refreshTokens:   refreshTokens,
		jwt:             jwt,
		hasher:          hasher,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// Register creates an account and immediately establishes a session.
func (s *authService) Register(ctx context.Context, in domain.RegisterInput) (*domain.AuthResult, error) {
	email := NormalizeEmail(in.Email)

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(in.Name),
		Email:        email,
		PasswordHash: hash,
		Role:         domain.RoleUser,
	}
	// The partial unique index on email is the authoritative duplicate check;
	// the repository translates violations into domain.ErrEmailTaken.
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

// Login verifies credentials and establishes a session.
func (s *authService) Login(ctx context.Context, in domain.LoginInput) (*domain.AuthResult, error) {
	user, err := s.users.GetByEmail(ctx, NormalizeEmail(in.Email))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// Equalize the failure path: hash anyway so response timing does
			// not reveal whether the email exists.
			_, _ = s.hasher.Hash(in.Password)
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !s.hasher.Verify(user.PasswordHash, in.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	return s.issueSession(ctx, user)
}

// Refresh rotates a refresh token into a new session (token pair).
func (s *authService) Refresh(ctx context.Context, refreshToken string) (*domain.AuthResult, error) {
	if refreshToken == "" {
		return nil, domain.ErrTokenInvalid
	}

	stored, err := s.refreshTokens.GetByHash(ctx, token.HashOpaque(refreshToken))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrTokenInvalid
		}
		return nil, err
	}

	now := time.Now().UTC()

	// Reuse detection: a revoked token being presented again is a strong
	// signal the token was stolen (rotation already consumed it once).
	// Kill every session of the user.
	if stored.Revoked() {
		slog.Warn("refresh token reuse detected; revoking all user sessions",
			"user_id", stored.UserID)
		if err := s.refreshTokens.RevokeAllForUser(ctx, stored.UserID); err != nil {
			return nil, err
		}
		return nil, domain.ErrTokenInvalid
	}

	if now.After(stored.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	user, err := s.users.GetByID(ctx, stored.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrTokenInvalid
		}
		return nil, err
	}

	// Rotation: consume the presented token, then issue a fresh pair.
	if err := s.refreshTokens.Revoke(ctx, stored.ID); err != nil {
		return nil, err
	}

	return s.issueSession(ctx, user)
}

// Logout revokes the presented refresh token. Unknown or empty tokens are
// ignored: logout must be idempotent and never fail the client.
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	stored, err := s.refreshTokens.GetByHash(ctx, token.HashOpaque(refreshToken))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}
	return s.refreshTokens.Revoke(ctx, stored.ID)
}

// EnsureAdmin seeds the bootstrap admin account (called on startup when
// ADMIN_EMAIL/ADMIN_PASSWORD are configured). Existing accounts are left
// untouched: this must never overwrite a password.
func (s *authService) EnsureAdmin(ctx context.Context, name, email, pass string) error {
	email = NormalizeEmail(email)

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil // account exists; nothing to do
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	hash, err := s.hasher.Hash(pass)
	if err != nil {
		return fmt.Errorf("ensure admin: %w", err)
	}
	admin := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         domain.RoleAdmin,
	}
	if err := s.users.Create(ctx, admin); err != nil {
		return err
	}
	slog.Info("bootstrap admin account created", "email", email)
	return nil
}

// issueSession generates the JWT access token and a rotating refresh token,
// persisting only the refresh token's hash.
func (s *authService) issueSession(ctx context.Context, user *domain.User) (*domain.AuthResult, error) {
	accessToken, accessExpiresAt, err := s.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.NewOpaque()
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().UTC().Add(s.refreshTokenTTL)
	record := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: token.HashOpaque(refreshToken),
		ExpiresAt: refreshExpiresAt,
	}
	if err := s.refreshTokens.Create(ctx, record); err != nil {
		return nil, err
	}

	return &domain.AuthResult{
		User:                  user,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

// NormalizeEmail canonicalizes an email address for storage and lookup.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

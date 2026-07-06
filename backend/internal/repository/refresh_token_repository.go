package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

// refreshTokenRepository is the GORM implementation of domain.RefreshTokenRepository.
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository builds the PostgreSQL-backed refresh token repository.
func NewRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("creating refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := r.db.WithContext(ctx).First(&token, "token_hash = ?", tokenHash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting refresh token: %w", err)
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", time.Now().UTC()).Error
	if err != nil {
		return fmt.Errorf("revoking refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Model(&domain.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", time.Now().UTC()).Error
	if err != nil {
		return fmt.Errorf("revoking user refresh tokens: %w", err)
	}
	return nil
}

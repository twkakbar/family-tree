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

type passwordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository builds the PostgreSQL-backed reset-token repo.
func NewPasswordResetTokenRepository(db *gorm.DB) domain.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("creating password reset token: %w", err)
	}
	return nil
}

func (r *passwordResetTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	var token domain.PasswordResetToken
	err := r.db.WithContext(ctx).First(&token, "token_hash = ?", tokenHash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting password reset token: %w", err)
	}
	return &token, nil
}

func (r *passwordResetTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&domain.PasswordResetToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now)
	if res.Error != nil {
		return fmt.Errorf("revoking password reset token: %w", res.Error)
	}
	return nil
}

func (r *passwordResetTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&domain.PasswordResetToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)
	if res.Error != nil {
		return fmt.Errorf("revoking password reset tokens for user: %w", res.Error)
	}
	return nil
}

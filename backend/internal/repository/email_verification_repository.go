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

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

// NewEmailVerificationTokenRepository builds the verification-token repo.
func NewEmailVerificationTokenRepository(db *gorm.DB) domain.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(ctx context.Context, token *domain.EmailVerificationToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("creating email verification token: %w", err)
	}
	return nil
}

func (r *emailVerificationTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.EmailVerificationToken, error) {
	var token domain.EmailVerificationToken
	err := r.db.WithContext(ctx).First(&token, "token_hash = ?", tokenHash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting email verification token: %w", err)
	}
	return &token, nil
}

func (r *emailVerificationTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&domain.EmailVerificationToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now)
	if res.Error != nil {
		return fmt.Errorf("revoking email verification token: %w", res.Error)
	}
	return nil
}

func (r *emailVerificationTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&domain.EmailVerificationToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)
	if res.Error != nil {
		return fmt.Errorf("revoking email verification tokens: %w", res.Error)
	}
	return nil
}

type emailDispatchRepository struct {
	db *gorm.DB
}

// NewEmailDispatchRepository builds the outbound-mail log.
func NewEmailDispatchRepository(db *gorm.DB) domain.EmailDispatchRepository {
	return &emailDispatchRepository{db: db}
}

func (r *emailDispatchRepository) Create(ctx context.Context, row *domain.EmailDispatch) error {
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("recording email dispatch: %w", err)
	}
	return nil
}

func (r *emailDispatchRepository) CountSince(ctx context.Context, email, purpose string, since time.Time) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&domain.EmailDispatch{}).
		Where("email = ? AND purpose = ? AND sent_at >= ?", email, purpose, since).
		Count(&n).Error
	if err != nil {
		return 0, fmt.Errorf("counting email dispatches: %w", err)
	}
	return n, nil
}

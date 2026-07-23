package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type editAccessRequestRepository struct {
	db *gorm.DB
}

// NewEditAccessRequestRepository builds the PostgreSQL-backed access request repository.
func NewEditAccessRequestRepository(db *gorm.DB) domain.EditAccessRequestRepository {
	return &editAccessRequestRepository{db: db}
}

func (r *editAccessRequestRepository) Create(ctx context.Context, req *domain.EditAccessRequest) error {
	if err := r.db.WithContext(ctx).Create(req).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAccessRequestPending
		}
		return fmt.Errorf("creating edit access request: %w", err)
	}
	return nil
}

func (r *editAccessRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.EditAccessRequest, error) {
	var req domain.EditAccessRequest
	err := r.db.WithContext(ctx).First(&req, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting edit access request: %w", err)
	}
	return &req, nil
}

func (r *editAccessRequestRepository) GetLatestForUser(ctx context.Context, userID uuid.UUID) (*domain.EditAccessRequest, error) {
	var req domain.EditAccessRequest
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting latest edit access request: %w", err)
	}
	return &req, nil
}

func (r *editAccessRequestRepository) ListByStatus(ctx context.Context, status domain.AccessRequestStatus) ([]domain.EditAccessRequest, error) {
	var reqs []domain.EditAccessRequest
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at ASC").
		Find(&reqs).Error
	if err != nil {
		return nil, fmt.Errorf("listing edit access requests: %w", err)
	}
	return reqs, nil
}

func (r *editAccessRequestRepository) Update(ctx context.Context, req *domain.EditAccessRequest) error {
	if err := r.db.WithContext(ctx).Save(req).Error; err != nil {
		return fmt.Errorf("updating edit access request: %w", err)
	}
	return nil
}

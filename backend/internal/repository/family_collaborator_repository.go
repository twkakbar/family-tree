package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type familyCollaboratorRepository struct {
	db *gorm.DB
}

// NewFamilyCollaboratorRepository builds the PostgreSQL-backed collaborator repository.
func NewFamilyCollaboratorRepository(db *gorm.DB) domain.FamilyCollaboratorRepository {
	return &familyCollaboratorRepository{db: db}
}

func (r *familyCollaboratorRepository) Create(ctx context.Context, collab *domain.FamilyCollaborator) error {
	if err := r.db.WithContext(ctx).Create(collab).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyCollaborator
		}
		return fmt.Errorf("creating family collaborator: %w", err)
	}
	return nil
}

func (r *familyCollaboratorRepository) Get(ctx context.Context, familyID, userID uuid.UUID) (*domain.FamilyCollaborator, error) {
	var collab domain.FamilyCollaborator
	err := r.db.WithContext(ctx).
		Where("family_id = ? AND user_id = ?", familyID, userID).
		First(&collab).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting family collaborator: %w", err)
	}
	return &collab, nil
}

func (r *familyCollaboratorRepository) ListByFamily(ctx context.Context, familyID uuid.UUID) ([]domain.FamilyCollaborator, error) {
	var collabs []domain.FamilyCollaborator
	err := r.db.WithContext(ctx).
		Where("family_id = ?", familyID).
		Order("created_at ASC").
		Find(&collabs).Error
	if err != nil {
		return nil, fmt.Errorf("listing family collaborators: %w", err)
	}
	return collabs, nil
}

func (r *familyCollaboratorRepository) Update(ctx context.Context, collab *domain.FamilyCollaborator) error {
	if err := r.db.WithContext(ctx).Save(collab).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return domain.ErrForbidden
		}
		return fmt.Errorf("updating family collaborator: %w", err)
	}
	return nil
}

func (r *familyCollaboratorRepository) Delete(ctx context.Context, familyID, userID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("family_id = ? AND user_id = ?", familyID, userID).
		Delete(&domain.FamilyCollaborator{})
	if res.Error != nil {
		return fmt.Errorf("deleting family collaborator: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *familyCollaboratorRepository) CountByRole(ctx context.Context, familyID uuid.UUID, role domain.FamilyRole) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.FamilyCollaborator{}).
		Where("family_id = ? AND role = ?", familyID, role).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting family collaborators by role: %w", err)
	}
	return count, nil
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type familyRepository struct {
	db *gorm.DB
}

// NewFamilyRepository builds the PostgreSQL-backed family repository.
func NewFamilyRepository(db *gorm.DB) domain.FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) Create(ctx context.Context, family *domain.Family) error {
	if err := r.db.WithContext(ctx).Create(family).Error; err != nil {
		return fmt.Errorf("creating family: %w", err)
	}
	return nil
}

func (r *familyRepository) CreateWithOwner(ctx context.Context, family *domain.Family, ownerID uuid.UUID, role domain.FamilyRole) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(family).Error; err != nil {
			return fmt.Errorf("creating family: %w", err)
		}
		collab := &domain.FamilyCollaborator{
			ID:       uuid.New(),
			FamilyID: family.ID,
			UserID:   ownerID,
			Role:     role,
		}
		if err := tx.Create(collab).Error; err != nil {
			return fmt.Errorf("creating owner collaborator: %w", err)
		}
		return nil
	})
}

func (r *familyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Family, error) {
	var family domain.Family
	err := r.db.WithContext(ctx).First(&family, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting family by id: %w", err)
	}
	return &family, nil
}

func (r *familyRepository) GetByName(ctx context.Context, name string) (*domain.Family, error) {
	var family domain.Family
	err := r.db.WithContext(ctx).First(&family, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting family by name: %w", err)
	}
	return &family, nil
}

func (r *familyRepository) Update(ctx context.Context, family *domain.Family) error {
	if err := r.db.WithContext(ctx).Save(family).Error; err != nil {
		return fmt.Errorf("updating family: %w", err)
	}
	return nil
}

func (r *familyRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&domain.Family{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("soft deleting family: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

type familySummaryRow struct {
	domain.Family
	Role         domain.FamilyRole `gorm:"column:role"`
	PersonsCount int64             `gorm:"column:persons_count"`
}

func (r *familyRepository) summaryQuery(userID uuid.UUID) *gorm.DB {
	return r.db.Table("families").
		Select(`
			families.*,
			family_collaborators.role AS role,
			(SELECT COUNT(*) FROM persons
			 WHERE persons.family_id = families.id AND persons.deleted_at IS NULL) AS persons_count`).
		Joins("INNER JOIN family_collaborators ON family_collaborators.family_id = families.id").
		Where("family_collaborators.user_id = ?", userID).
		Where("families.deleted_at IS NULL")
}

func (r *familyRepository) ListForUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]domain.FamilySummary, int64, error) {
	base := r.summaryQuery(userID)

	var total int64
	if err := base.WithContext(ctx).Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting families for user: %w", err)
	}

	var rows []familySummaryRow
	err := base.WithContext(ctx).
		Order("families.created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listing families for user: %w", err)
	}

	out := make([]domain.FamilySummary, len(rows))
	for i, row := range rows {
		out[i] = domain.FamilySummary{
			Family:       row.Family,
			Role:         row.Role,
			PersonsCount: row.PersonsCount,
		}
	}
	return out, total, nil
}

func (r *familyRepository) GetSummaryForUser(ctx context.Context, userID, familyID uuid.UUID) (*domain.FamilySummary, error) {
	var row familySummaryRow
	err := r.summaryQuery(userID).
		WithContext(ctx).
		Where("families.id = ?", familyID).
		Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("getting family summary: %w", err)
	}
	if row.ID == uuid.Nil {
		return nil, domain.ErrNotFound
	}
	return &domain.FamilySummary{
		Family:       row.Family,
		Role:         row.Role,
		PersonsCount: row.PersonsCount,
	}, nil
}

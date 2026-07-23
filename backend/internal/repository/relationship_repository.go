package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type relationshipRepository struct {
	db *gorm.DB
}

// NewRelationshipRepository builds the PostgreSQL-backed relationship repository.
func NewRelationshipRepository(db *gorm.DB) domain.RelationshipRepository {
	return &relationshipRepository{db: db}
}

func (r *relationshipRepository) Create(ctx context.Context, rel *domain.Relationship) error {
	if err := r.db.WithContext(ctx).Create(rel).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrRelationshipExists
		}
		return fmt.Errorf("creating relationship: %w", err)
	}
	return nil
}

func (r *relationshipRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Relationship, error) {
	var rel domain.Relationship
	err := r.db.WithContext(ctx).First(&rel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting relationship by id: %w", err)
	}
	return &rel, nil
}

func (r *relationshipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&domain.Relationship{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("deleting relationship: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *relationshipRepository) ListByFamily(ctx context.Context, familyID uuid.UUID) ([]domain.Relationship, error) {
	var rels []domain.Relationship
	err := r.db.WithContext(ctx).
		Where("family_id = ?", familyID).
		Order("created_at ASC").
		Find(&rels).Error
	if err != nil {
		return nil, fmt.Errorf("listing relationships by family: %w", err)
	}
	return rels, nil
}

func (r *relationshipRepository) ListForPerson(ctx context.Context, personID uuid.UUID) ([]domain.Relationship, error) {
	var rels []domain.Relationship
	err := r.db.WithContext(ctx).
		Where("from_person_id = ? OR to_person_id = ?", personID, personID).
		Order("created_at ASC").
		Find(&rels).Error
	if err != nil {
		return nil, fmt.Errorf("listing relationships for person: %w", err)
	}
	return rels, nil
}

func (r *relationshipRepository) CountParentsOf(ctx context.Context, childID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Relationship{}).
		Where("to_person_id = ? AND type = ?", childID, domain.RelationshipParent).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting parents: %w", err)
	}
	return count, nil
}

const ancestorPathSQL = `
WITH RECURSIVE lineage AS (
    SELECT r.from_person_id AS person_id
    FROM relationships r
    WHERE r.to_person_id = ? AND r.type = 'parent'
  UNION ALL
    SELECT r.from_person_id
    FROM relationships r
    JOIN lineage l ON r.to_person_id = l.person_id
    WHERE r.type = 'parent'
)
SELECT EXISTS (
    SELECT 1 FROM lineage WHERE person_id = ?
) AS found`

func (r *relationshipRepository) HasAncestorPath(ctx context.Context, startID, targetID uuid.UUID) (bool, error) {
	var found bool
	err := r.db.WithContext(ctx).Raw(ancestorPathSQL, startID, targetID).Scan(&found).Error
	if err != nil {
		return false, fmt.Errorf("checking ancestor path: %w", err)
	}
	return found, nil
}

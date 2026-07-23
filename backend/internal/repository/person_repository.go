package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type personRepository struct {
	db *gorm.DB
}

// NewPersonRepository builds the PostgreSQL-backed person repository.
func NewPersonRepository(db *gorm.DB) domain.PersonRepository {
	return &personRepository{db: db}
}

func (r *personRepository) CountByFamily(ctx context.Context, familyID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Person{}).
		Where("family_id = ?", familyID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("counting persons by family: %w", err)
	}
	return count, nil
}

func (r *personRepository) Create(ctx context.Context, person *domain.Person) error {
	if err := r.db.WithContext(ctx).Create(person).Error; err != nil {
		return fmt.Errorf("creating person: %w", err)
	}
	return nil
}

func (r *personRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error) {
	var person domain.Person
	err := r.db.WithContext(ctx).First(&person, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("getting person by id: %w", err)
	}
	return &person, nil
}

func (r *personRepository) Update(ctx context.Context, person *domain.Person) error {
	if err := r.db.WithContext(ctx).Save(person).Error; err != nil {
		return fmt.Errorf("updating person: %w", err)
	}
	return nil
}

func (r *personRepository) SoftDeleteWithRelationships(ctx context.Context, personID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Where("from_person_id = ? OR to_person_id = ?", personID, personID).
			Delete(&domain.Relationship{})
		if res.Error != nil {
			return fmt.Errorf("deleting person relationships: %w", res.Error)
		}

		res = tx.Delete(&domain.Person{}, "id = ?", personID)
		if res.Error != nil {
			return fmt.Errorf("soft deleting person: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		return nil
	})
}

func (r *personRepository) listQuery(familyID uuid.UUID, search string) *gorm.DB {
	q := r.db.Model(&domain.Person{}).Where("family_id = ?", familyID)
	if search != "" {
		// Trigram-backed pattern match on the full display name (index:
		// persons_name_trgm_idx). Also match either name column for short queries.
		pattern := "%" + search + "%"
		q = q.Where(
			"(first_name || ' ' || last_name) ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			pattern, pattern, pattern,
		)
	}
	return q
}

func (r *personRepository) ListByFamily(ctx context.Context, familyID uuid.UUID, search string, offset, limit int) ([]domain.Person, int64, error) {
	base := r.listQuery(familyID, search)

	var total int64
	if err := base.WithContext(ctx).Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting persons: %w", err)
	}

	var persons []domain.Person
	q := base.WithContext(ctx)
	if search != "" {
		// Rank closer name matches first when searching.
		q = q.Order(gorm.Expr(
			"similarity(first_name || ' ' || last_name, ?) DESC, last_name ASC, first_name ASC",
			search,
		))
	} else {
		q = q.Order("last_name ASC, first_name ASC")
	}
	err := q.Offset(offset).Limit(limit).Find(&persons).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listing persons: %w", err)
	}
	return persons, total, nil
}

func (r *personRepository) ListByFamilyAll(ctx context.Context, familyID uuid.UUID) ([]domain.Person, error) {
	var persons []domain.Person
	err := r.db.WithContext(ctx).
		Where("family_id = ?", familyID).
		Order("last_name ASC, first_name ASC").
		Find(&persons).Error
	if err != nil {
		return nil, fmt.Errorf("listing all persons: %w", err)
	}
	return persons, nil
}

type personDepthRow struct {
	domain.Person
	Depth int `gorm:"column:depth"`
}

func (r *personRepository) ListAncestors(ctx context.Context, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	return r.traverse(ctx, personID, depth, ancestorsSQL)
}

func (r *personRepository) ListDescendants(ctx context.Context, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	return r.traverse(ctx, personID, depth, descendantsSQL)
}

const ancestorsSQL = `
WITH RECURSIVE lineage AS (
    SELECT r.from_person_id AS person_id, 1 AS depth
    FROM relationships r
    JOIN persons child ON child.id = r.to_person_id AND child.deleted_at IS NULL
    WHERE r.to_person_id = ? AND r.type = 'parent'
  UNION ALL
    SELECT r.from_person_id, l.depth + 1
    FROM relationships r
    JOIN lineage l ON r.to_person_id = l.person_id
    WHERE r.type = 'parent' AND l.depth < ?
)
SELECT p.*, l.depth
FROM lineage l
JOIN persons p ON p.id = l.person_id AND p.deleted_at IS NULL
ORDER BY l.depth ASC, p.last_name ASC, p.first_name ASC`

const descendantsSQL = `
WITH RECURSIVE lineage AS (
    SELECT r.to_person_id AS person_id, 1 AS depth
    FROM relationships r
    JOIN persons parent ON parent.id = r.from_person_id AND parent.deleted_at IS NULL
    WHERE r.from_person_id = ? AND r.type = 'parent'
  UNION ALL
    SELECT r.to_person_id, l.depth + 1
    FROM relationships r
    JOIN lineage l ON r.from_person_id = l.person_id
    WHERE r.type = 'parent' AND l.depth < ?
)
SELECT p.*, l.depth
FROM lineage l
JOIN persons p ON p.id = l.person_id AND p.deleted_at IS NULL
ORDER BY l.depth ASC, p.last_name ASC, p.first_name ASC`

func (r *personRepository) traverse(ctx context.Context, personID uuid.UUID, depth int, sql string) ([]domain.PersonWithDepth, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 20 {
		depth = 20
	}

	var rows []personDepthRow
	if err := r.db.WithContext(ctx).Raw(sql, personID, depth).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("traversing persons: %w", err)
	}

	out := make([]domain.PersonWithDepth, len(rows))
	for i, row := range rows {
		out[i] = domain.PersonWithDepth{Person: row.Person, Depth: row.Depth}
	}
	return out, nil
}

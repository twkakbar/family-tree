package service

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/upload"
)

type personService struct {
	families      domain.FamilyRepository
	persons       domain.PersonRepository
	relationships domain.RelationshipRepository
	users         domain.UserRepository
	uploadDir     string
	maxUpload     int64
}

// NewPersonService wires person CRUD, tree data, relationships, and photo uploads.
func NewPersonService(
	families domain.FamilyRepository,
	persons domain.PersonRepository,
	relationships domain.RelationshipRepository,
	users domain.UserRepository,
	uploadDir string,
	maxUploadBytes int64,
) domain.PersonService {
	if maxUploadBytes <= 0 {
		maxUploadBytes = 10 << 20
	}
	return &personService{
		families:      families,
		persons:       persons,
		relationships: relationships,
		users:         users,
		uploadDir:     uploadDir,
		maxUpload:     maxUploadBytes,
	}
}

func (s *personService) ListByFamily(ctx context.Context, userID, familyID uuid.UUID, search string, page, perPage int) ([]domain.Person, int64, error) {
	if _, err := s.requireFamily(ctx, userID, familyID, domain.FamilyRoleViewer); err != nil {
		return nil, 0, err
	}
	params := normalizePage(page, perPage)
	return s.persons.ListByFamily(ctx, familyID, strings.TrimSpace(search), params.offset(), params.limit())
}

func (s *personService) Create(ctx context.Context, userID, familyID uuid.UUID, in domain.CreatePersonInput) (*domain.Person, error) {
	if _, err := s.requireFamily(ctx, userID, familyID, domain.FamilyRoleEditor); err != nil {
		return nil, err
	}
	if !in.Gender.Valid() {
		return nil, domain.ErrInvalidRelationship
	}
	if err := validateDates(in.BirthDate, in.DeathDate); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	person := &domain.Person{
		ID:        uuid.New(),
		FamilyID:  familyID,
		FirstName: strings.TrimSpace(in.FirstName),
		LastName:  strings.TrimSpace(in.LastName),
		Gender:    in.Gender,
		BirthDate: in.BirthDate,
		DeathDate: in.DeathDate,
		Bio:       strings.TrimSpace(in.Bio),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.persons.Create(ctx, person); err != nil {
		return nil, err
	}
	return person, nil
}

func (s *personService) GetByID(ctx context.Context, userID, personID uuid.UUID) (*domain.PersonDetail, error) {
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireFamily(ctx, userID, person.FamilyID, domain.FamilyRoleViewer); err != nil {
		return nil, err
	}
	rels, err := s.relationships.ListForPerson(ctx, personID)
	if err != nil {
		return nil, err
	}
	return &domain.PersonDetail{Person: *person, Relationships: rels}, nil
}

func (s *personService) Update(ctx context.Context, userID, personID uuid.UUID, in domain.UpdatePersonInput) (*domain.Person, error) {
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireFamily(ctx, userID, person.FamilyID, domain.FamilyRoleEditor); err != nil {
		return nil, err
	}
	if !in.Gender.Valid() {
		return nil, domain.ErrInvalidRelationship
	}
	if err := validateDates(in.BirthDate, in.DeathDate); err != nil {
		return nil, err
	}

	person.FirstName = strings.TrimSpace(in.FirstName)
	person.LastName = strings.TrimSpace(in.LastName)
	person.Gender = in.Gender
	person.BirthDate = in.BirthDate
	person.DeathDate = in.DeathDate
	person.Bio = strings.TrimSpace(in.Bio)
	person.UpdatedAt = time.Now().UTC()
	if err := s.persons.Update(ctx, person); err != nil {
		return nil, err
	}
	return person, nil
}

func (s *personService) Delete(ctx context.Context, userID, personID uuid.UUID) error {
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return err
	}
	if _, err := s.requireFamily(ctx, userID, person.FamilyID, domain.FamilyRoleEditor); err != nil {
		return err
	}
	return s.persons.SoftDeleteWithRelationships(ctx, personID)
}

func (s *personService) GetFamilyTree(ctx context.Context, userID, familyID uuid.UUID) (*domain.FamilyTree, error) {
	if _, err := s.requireFamily(ctx, userID, familyID, domain.FamilyRoleViewer); err != nil {
		return nil, err
	}
	persons, err := s.persons.ListByFamilyAll(ctx, familyID)
	if err != nil {
		return nil, err
	}
	rels, err := s.relationships.ListByFamily(ctx, familyID)
	if err != nil {
		return nil, err
	}
	return &domain.FamilyTree{Persons: persons, Relationships: rels}, nil
}

func (s *personService) GetAncestors(ctx context.Context, userID, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireFamily(ctx, userID, person.FamilyID, domain.FamilyRoleViewer); err != nil {
		return nil, err
	}
	return s.persons.ListAncestors(ctx, personID, normalizeDepth(depth))
}

func (s *personService) GetDescendants(ctx context.Context, userID, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireFamily(ctx, userID, person.FamilyID, domain.FamilyRoleViewer); err != nil {
		return nil, err
	}
	return s.persons.ListDescendants(ctx, personID, normalizeDepth(depth))
}

func (s *personService) CreateRelationship(ctx context.Context, userID, familyID uuid.UUID, in domain.CreateRelationshipInput) (*domain.Relationship, error) {
	if _, err := s.requireFamily(ctx, userID, familyID, domain.FamilyRoleEditor); err != nil {
		return nil, err
	}

	fromPerson, err := s.persons.GetByID(ctx, in.FromPersonID)
	if err != nil {
		return nil, err
	}
	toPerson, err := s.persons.GetByID(ctx, in.ToPersonID)
	if err != nil {
		return nil, err
	}
	if fromPerson.FamilyID != familyID || toPerson.FamilyID != familyID {
		return nil, domain.ErrInvalidRelationship
	}
	if in.FromPersonID == in.ToPersonID {
		return nil, domain.ErrInvalidRelationship
	}

	fromID := in.FromPersonID
	toID := in.ToPersonID
	relType := in.Type

	switch relType {
	case domain.RelationshipParent:
		parentCount, err := s.relationships.CountParentsOf(ctx, toID)
		if err != nil {
			return nil, err
		}
		if parentCount >= 2 {
			return nil, domain.ErrInvalidRelationship
		}
		cycle, err := s.relationships.HasAncestorPath(ctx, fromID, toID)
		if err != nil {
			return nil, err
		}
		if cycle {
			return nil, domain.ErrInvalidRelationship
		}
	case domain.RelationshipSpouse:
		fromID, toID = canonicalSpousePair(fromID, toID)
	default:
		return nil, domain.ErrInvalidRelationship
	}

	now := time.Now().UTC()
	rel := &domain.Relationship{
		ID:           uuid.New(),
		FamilyID:     familyID,
		FromPersonID: fromID,
		ToPersonID:   toID,
		Type:         relType,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.relationships.Create(ctx, rel); err != nil {
		return nil, err
	}
	return rel, nil
}

func (s *personService) DeleteRelationship(ctx context.Context, userID, relationshipID uuid.UUID) error {
	rel, err := s.relationships.GetByID(ctx, relationshipID)
	if err != nil {
		return err
	}
	if _, err := s.requireFamily(ctx, userID, rel.FamilyID, domain.FamilyRoleEditor); err != nil {
		return err
	}
	return s.relationships.Delete(ctx, relationshipID)
}

func (s *personService) UpdatePhoto(ctx context.Context, userID, personID uuid.UUID, data []byte) (*domain.Person, error) {
	if len(data) == 0 {
		return nil, domain.ErrInvalidUpload
	}
	if int64(len(data)) > s.maxUpload {
		return nil, domain.ErrInvalidUpload
	}

	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireFamily(ctx, userID, person.FamilyID, domain.FamilyRoleEditor); err != nil {
		return nil, err
	}

	_, ext, err := upload.DetectImageMIME(data)
	if err != nil {
		return nil, domain.ErrInvalidUpload
	}

	relative := filepath.ToSlash(filepath.Join("persons", person.ID.String()+ext))
	if err := upload.Save(s.uploadDir, relative, data); err != nil {
		return nil, fmt.Errorf("saving person photo: %w", err)
	}

	old := person.PhotoPath
	person.PhotoPath = relative
	person.UpdatedAt = time.Now().UTC()
	if err := s.persons.Update(ctx, person); err != nil {
		_ = upload.Remove(s.uploadDir, relative)
		return nil, err
	}
	if old != "" && old != relative {
		_ = upload.Remove(s.uploadDir, old)
	}
	return person, nil
}

func (s *personService) requireFamily(ctx context.Context, userID, familyID uuid.UUID, minRole domain.FamilyRole) (*domain.FamilySummary, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.IsAdmin() {
		family, err := s.families.GetByID(ctx, familyID)
		if err != nil {
			return nil, err
		}
		count, err := s.persons.CountByFamily(ctx, familyID)
		if err != nil {
			return nil, err
		}
		return &domain.FamilySummary{
			Family:       *family,
			Role:         domain.FamilyRoleOwner,
			PersonsCount: count,
		}, nil
	}

	summary, err := s.families.GetSummaryForUser(ctx, userID, familyID)
	if err != nil {
		return nil, err
	}
	if !summary.Role.AtLeast(minRole) {
		return nil, domain.ErrForbidden
	}
	return summary, nil
}

func validateDates(birth, death *time.Time) error {
	if birth != nil && death != nil && death.Before(*birth) {
		return domain.ErrInvalidRelationship
	}
	return nil
}

func canonicalSpousePair(a, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if bytes.Compare(a[:], b[:]) > 0 {
		return b, a
	}
	return a, b
}

func normalizeDepth(depth int) int {
	if depth < 1 {
		return 5
	}
	if depth > 20 {
		return 20
	}
	return depth
}

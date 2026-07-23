package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

const maxImportPersons = 2000

type transferService struct {
	families domain.FamilyService
	persons  domain.PersonService
}

// NewTransferService wires JSON tree export/import.
func NewTransferService(families domain.FamilyService, persons domain.PersonService) domain.TransferService {
	return &transferService{families: families, persons: persons}
}

func (s *transferService) ExportJSON(ctx context.Context, userID, familyID uuid.UUID) (*domain.TreeTransfer, error) {
	summary, err := s.families.GetByID(ctx, userID, familyID)
	if err != nil {
		return nil, err
	}
	tree, err := s.persons.GetFamilyTree(ctx, userID, familyID)
	if err != nil {
		return nil, err
	}

	out := &domain.TreeTransfer{
		FormatVersion: domain.TreeTransferFormatVersion,
		ExportedAt:    time.Now().UTC(),
		Family: domain.TreeTransferFamily{
			Name:        summary.Family.Name,
			Description: summary.Family.Description,
		},
		Persons:       make([]domain.TreeTransferPerson, 0, len(tree.Persons)),
		Relationships: make([]domain.TreeTransferEdge, 0, len(tree.Relationships)),
	}

	for i := range tree.Persons {
		p := &tree.Persons[i]
		tp := domain.TreeTransferPerson{
			ID:        p.ID.String(),
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Gender:    p.Gender,
			Bio:       p.Bio,
		}
		if p.BirthDate != nil {
			s := p.BirthDate.UTC().Format("2006-01-02")
			tp.BirthDate = &s
		}
		if p.DeathDate != nil {
			s := p.DeathDate.UTC().Format("2006-01-02")
			tp.DeathDate = &s
		}
		out.Persons = append(out.Persons, tp)
	}

	for i := range tree.Relationships {
		r := &tree.Relationships[i]
		out.Relationships = append(out.Relationships, domain.TreeTransferEdge{
			FromPersonID: r.FromPersonID.String(),
			ToPersonID:   r.ToPersonID.String(),
			Type:         r.Type,
		})
	}
	return out, nil
}

func (s *transferService) ImportJSON(ctx context.Context, userID, familyID uuid.UUID, doc *domain.TreeTransfer) (*domain.ImportResult, error) {
	if doc == nil || doc.FormatVersion != domain.TreeTransferFormatVersion {
		return nil, domain.ErrUnsupportedFormat
	}
	if len(doc.Persons) > maxImportPersons {
		return nil, domain.ErrUnsupportedFormat
	}

	// Editor+ gate via creating a no-op check: list tree requires viewer;
	// Create requires editor — probe with GetFamilyTree then Create path.
	if _, err := s.persons.GetFamilyTree(ctx, userID, familyID); err != nil {
		return nil, err
	}

	idMap := make(map[string]uuid.UUID, len(doc.Persons))
	createdPersons := 0

	for _, raw := range doc.Persons {
		oldID := strings.TrimSpace(raw.ID)
		if oldID == "" {
			return nil, domain.ErrUnsupportedFormat
		}
		if _, err := uuid.Parse(oldID); err != nil {
			return nil, domain.ErrUnsupportedFormat
		}
		gender := raw.Gender
		if !gender.Valid() {
			gender = domain.GenderUnknown
		}
		birth, err := parseOptionalDate(raw.BirthDate)
		if err != nil {
			return nil, domain.ErrUnsupportedFormat
		}
		death, err := parseOptionalDate(raw.DeathDate)
		if err != nil {
			return nil, domain.ErrUnsupportedFormat
		}
		person, err := s.persons.Create(ctx, userID, familyID, domain.CreatePersonInput{
			FirstName: strings.TrimSpace(raw.FirstName),
			LastName:  strings.TrimSpace(raw.LastName),
			Gender:    gender,
			BirthDate: birth,
			DeathDate: death,
			Bio:       strings.TrimSpace(raw.Bio),
		})
		if err != nil {
			return nil, fmt.Errorf("import person %s: %w", oldID, err)
		}
		idMap[oldID] = person.ID
		createdPersons++
	}

	createdRels := 0
	for _, edge := range doc.Relationships {
		fromID, okFrom := idMap[strings.TrimSpace(edge.FromPersonID)]
		toID, okTo := idMap[strings.TrimSpace(edge.ToPersonID)]
		if !okFrom || !okTo {
			return nil, domain.ErrUnsupportedFormat
		}
		_, err := s.persons.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
			FromPersonID: fromID,
			ToPersonID:   toID,
			Type:         edge.Type,
		})
		if err != nil {
			return nil, fmt.Errorf("import relationship: %w", err)
		}
		createdRels++
	}

	return &domain.ImportResult{
		PersonsCreated:       createdPersons,
		RelationshipsCreated: createdRels,
	}, nil
}

func parseOptionalDate(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*value))
	if err != nil {
		return nil, err
	}
	utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &utc, nil
}

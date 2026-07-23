package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type publicTreeService struct {
	royal   domain.RoyalFamilyService
	persons domain.PersonRepository
	rels    domain.RelationshipRepository
}

// NewPublicTreeService wires unauthenticated royal-tree reads.
func NewPublicTreeService(
	royal domain.RoyalFamilyService,
	persons domain.PersonRepository,
	rels domain.RelationshipRepository,
) domain.PublicTreeService {
	return &publicTreeService{royal: royal, persons: persons, rels: rels}
}

func (s *publicTreeService) Meta(ctx context.Context) (*domain.Family, error) {
	return s.royal.Get(ctx)
}

func (s *publicTreeService) Tree(ctx context.Context) (*domain.FamilyTree, error) {
	familyID, err := s.royal.ID(ctx)
	if err != nil {
		return nil, err
	}
	persons, err := s.persons.ListByFamilyAll(ctx, familyID)
	if err != nil {
		return nil, err
	}
	rels, err := s.rels.ListByFamily(ctx, familyID)
	if err != nil {
		return nil, err
	}
	return &domain.FamilyTree{Persons: persons, Relationships: rels}, nil
}

func (s *publicTreeService) GetPerson(ctx context.Context, personID uuid.UUID) (*domain.PersonDetail, error) {
	familyID, err := s.royal.ID(ctx)
	if err != nil {
		return nil, err
	}
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}
	if person.FamilyID != familyID {
		return nil, domain.ErrNotFound
	}
	rels, err := s.rels.ListForPerson(ctx, personID)
	if err != nil {
		return nil, err
	}
	return &domain.PersonDetail{Person: *person, Relationships: rels}, nil
}

func (s *publicTreeService) Ancestors(ctx context.Context, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	if err := s.ensureRoyalPerson(ctx, personID); err != nil {
		return nil, err
	}
	return s.persons.ListAncestors(ctx, personID, normalizeDepth(depth))
}

func (s *publicTreeService) Descendants(ctx context.Context, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	if err := s.ensureRoyalPerson(ctx, personID); err != nil {
		return nil, err
	}
	return s.persons.ListDescendants(ctx, personID, normalizeDepth(depth))
}

func (s *publicTreeService) ensureRoyalPerson(ctx context.Context, personID uuid.UUID) error {
	familyID, err := s.royal.ID(ctx)
	if err != nil {
		return err
	}
	person, err := s.persons.GetByID(ctx, personID)
	if err != nil {
		return err
	}
	if person.FamilyID != familyID {
		return domain.ErrNotFound
	}
	return nil
}

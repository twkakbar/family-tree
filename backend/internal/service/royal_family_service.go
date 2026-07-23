package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type royalFamilyService struct {
	families      domain.FamilyRepository
	collaborators domain.FamilyCollaboratorRepository
	users         domain.UserRepository
	name          string
}

// NewRoyalFamilyService resolves and seeds the singleton public family.
func NewRoyalFamilyService(
	families domain.FamilyRepository,
	collaborators domain.FamilyCollaboratorRepository,
	users domain.UserRepository,
	name string,
) domain.RoyalFamilyService {
	if strings.TrimSpace(name) == "" {
		name = "Kesultanan Aceh"
	}
	return &royalFamilyService{
		families:      families,
		collaborators: collaborators,
		users:         users,
		name:          strings.TrimSpace(name),
	}
}

func (s *royalFamilyService) Ensure(ctx context.Context, ownerID uuid.UUID, name, description string) (*domain.Family, error) {
	if strings.TrimSpace(name) != "" {
		s.name = strings.TrimSpace(name)
	}
	existing, err := s.families.GetByName(ctx, s.name)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	now := time.Now().UTC()
	family := &domain.Family{
		ID:          uuid.New(),
		Name:        s.name,
		Description: strings.TrimSpace(description),
		CreatedBy:   ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.families.CreateWithOwner(ctx, family, ownerID, domain.FamilyRoleOwner); err != nil {
		// Race: another instance created it.
		if existing, getErr := s.families.GetByName(ctx, s.name); getErr == nil {
			return existing, nil
		}
		return nil, err
	}
	slog.Info("royal family seeded", "id", family.ID, "name", family.Name)
	return family, nil
}

func (s *royalFamilyService) ID(ctx context.Context) (uuid.UUID, error) {
	family, err := s.Get(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return family.ID, nil
}

func (s *royalFamilyService) Get(ctx context.Context) (*domain.Family, error) {
	return s.families.GetByName(ctx, s.name)
}

func (s *royalFamilyService) CanEdit(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user.IsAdmin() {
		return true, nil
	}
	familyID, err := s.ID(ctx)
	if err != nil {
		return false, err
	}
	collab, err := s.collaborators.Get(ctx, familyID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return collab.Role.AtLeast(domain.FamilyRoleEditor), nil
}

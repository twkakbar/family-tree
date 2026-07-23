package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type familyService struct {
	families      domain.FamilyRepository
	collaborators domain.FamilyCollaboratorRepository
	users         domain.UserRepository
}

// NewFamilyService wires family CRUD and collaborator management.
func NewFamilyService(
	families domain.FamilyRepository,
	collaborators domain.FamilyCollaboratorRepository,
	users domain.UserRepository,
) domain.FamilyService {
	return &familyService{
		families:      families,
		collaborators: collaborators,
		users:         users,
	}
}

func (s *familyService) Create(ctx context.Context, userID uuid.UUID, in domain.CreateFamilyInput) (*domain.FamilySummary, error) {
	now := time.Now().UTC()
	family := &domain.Family{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(in.Name),
		Description: strings.TrimSpace(in.Description),
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.families.CreateWithOwner(ctx, family, userID, domain.FamilyRoleOwner); err != nil {
		return nil, err
	}
	return &domain.FamilySummary{
		Family:       *family,
		Role:         domain.FamilyRoleOwner,
		PersonsCount: 0,
	}, nil
}

func (s *familyService) ListForUser(ctx context.Context, userID uuid.UUID, page, perPage int) ([]domain.FamilySummary, int64, error) {
	params := normalizePage(page, perPage)
	return s.families.ListForUser(ctx, userID, params.offset(), params.limit())
}

func (s *familyService) GetByID(ctx context.Context, userID, familyID uuid.UUID) (*domain.FamilySummary, error) {
	return s.requireSummary(ctx, userID, familyID, domain.FamilyRoleViewer)
}

func (s *familyService) Update(ctx context.Context, userID, familyID uuid.UUID, in domain.UpdateFamilyInput) (*domain.FamilySummary, error) {
	summary, err := s.requireSummary(ctx, userID, familyID, domain.FamilyRoleOwner)
	if err != nil {
		return nil, err
	}

	summary.Family.Name = strings.TrimSpace(in.Name)
	summary.Family.Description = strings.TrimSpace(in.Description)
	summary.Family.UpdatedAt = time.Now().UTC()
	if err := s.families.Update(ctx, &summary.Family); err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *familyService) Delete(ctx context.Context, userID, familyID uuid.UUID) error {
	if _, err := s.requireSummary(ctx, userID, familyID, domain.FamilyRoleOwner); err != nil {
		return err
	}
	return s.families.SoftDelete(ctx, familyID)
}

func (s *familyService) ListCollaborators(ctx context.Context, userID, familyID uuid.UUID) ([]domain.CollaboratorSummary, error) {
	if _, err := s.requireSummary(ctx, userID, familyID, domain.FamilyRoleViewer); err != nil {
		return nil, err
	}
	return s.loadCollaboratorSummaries(ctx, familyID)
}

func (s *familyService) AddCollaborator(ctx context.Context, userID, familyID uuid.UUID, in domain.AddCollaboratorInput) (*domain.CollaboratorSummary, error) {
	if _, err := s.requireSummary(ctx, userID, familyID, domain.FamilyRoleOwner); err != nil {
		return nil, err
	}

	role := domain.FamilyRole(strings.TrimSpace(string(in.Role)))
	if role != domain.FamilyRoleEditor && role != domain.FamilyRoleViewer {
		return nil, domain.ErrForbidden
	}

	invitee, err := s.users.GetByEmail(ctx, strings.TrimSpace(in.Email))
	if err != nil {
		return nil, err
	}

	collab := &domain.FamilyCollaborator{
		ID:       uuid.New(),
		FamilyID: familyID,
		UserID:   invitee.ID,
		Role:     role,
	}
	if err := s.collaborators.Create(ctx, collab); err != nil {
		return nil, err
	}
	return &domain.CollaboratorSummary{
		Collaborator: *collab,
		User:         *invitee,
	}, nil
}

func (s *familyService) UpdateCollaboratorRole(ctx context.Context, userID, familyID, targetUserID uuid.UUID, role domain.FamilyRole) (*domain.CollaboratorSummary, error) {
	if _, err := s.requireSummary(ctx, userID, familyID, domain.FamilyRoleOwner); err != nil {
		return nil, err
	}
	if !role.Valid() {
		return nil, domain.ErrForbidden
	}

	collab, err := s.collaborators.Get(ctx, familyID, targetUserID)
	if err != nil {
		return nil, err
	}
	if collab.Role == domain.FamilyRoleOwner && role != domain.FamilyRoleOwner {
		if err := s.ensureNotLastOwner(ctx, familyID); err != nil {
			return nil, err
		}
	}

	collab.Role = role
	collab.UpdatedAt = time.Now().UTC()
	if err := s.collaborators.Update(ctx, collab); err != nil {
		return nil, err
	}

	user, err := s.users.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	return &domain.CollaboratorSummary{Collaborator: *collab, User: *user}, nil
}

func (s *familyService) RemoveCollaborator(ctx context.Context, userID, familyID, targetUserID uuid.UUID) error {
	if _, err := s.requireSummary(ctx, userID, familyID, domain.FamilyRoleOwner); err != nil {
		return err
	}

	collab, err := s.collaborators.Get(ctx, familyID, targetUserID)
	if err != nil {
		return err
	}
	if collab.Role == domain.FamilyRoleOwner {
		if err := s.ensureNotLastOwner(ctx, familyID); err != nil {
			return err
		}
	}
	return s.collaborators.Delete(ctx, familyID, targetUserID)
}

func (s *familyService) requireSummary(ctx context.Context, userID, familyID uuid.UUID, minRole domain.FamilyRole) (*domain.FamilySummary, error) {
	summary, err := s.families.GetSummaryForUser(ctx, userID, familyID)
	if err != nil {
		return nil, err
	}
	if !summary.Role.AtLeast(minRole) {
		return nil, domain.ErrForbidden
	}
	return summary, nil
}

func (s *familyService) ensureNotLastOwner(ctx context.Context, familyID uuid.UUID) error {
	count, err := s.collaborators.CountByRole(ctx, familyID, domain.FamilyRoleOwner)
	if err != nil {
		return err
	}
	if count <= 1 {
		return domain.ErrLastOwner
	}
	return nil
}

func (s *familyService) loadCollaboratorSummaries(ctx context.Context, familyID uuid.UUID) ([]domain.CollaboratorSummary, error) {
	collabs, err := s.collaborators.ListByFamily(ctx, familyID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CollaboratorSummary, 0, len(collabs))
	for _, collab := range collabs {
		user, err := s.users.GetByID(ctx, collab.UserID)
		if err != nil {
			return nil, err
		}
		out = append(out, domain.CollaboratorSummary{
			Collaborator: collab,
			User:         *user,
		})
	}
	return out, nil
}

type pageParams struct {
	page    int
	perPage int
}

func normalizePage(page, perPage int) pageParams {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return pageParams{page: page, perPage: perPage}
}

func (p pageParams) offset() int { return (p.page - 1) * p.perPage }
func (p pageParams) limit() int  { return p.perPage }

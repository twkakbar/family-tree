package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type accessRequestService struct {
	requests      domain.EditAccessRequestRepository
	users         domain.UserRepository
	royal         domain.RoyalFamilyService
	collaborators domain.FamilyCollaboratorRepository
}

// NewAccessRequestService wires edit-access request flow.
func NewAccessRequestService(
	requests domain.EditAccessRequestRepository,
	users domain.UserRepository,
	royal domain.RoyalFamilyService,
	collaborators domain.FamilyCollaboratorRepository,
) domain.AccessRequestService {
	return &accessRequestService{
		requests:      requests,
		users:         users,
		royal:         royal,
		collaborators: collaborators,
	}
}

func (s *accessRequestService) Submit(ctx context.Context, userID uuid.UUID, message string) (*domain.EditAccessRequest, error) {
	canEdit, err := s.royal.CanEdit(ctx, userID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if canEdit {
		return nil, domain.ErrAlreadyEditor
	}

	latest, err := s.requests.GetLatestForUser(ctx, userID)
	if err == nil && latest.Status == domain.AccessRequestPending {
		return nil, domain.ErrAccessRequestPending
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	now := time.Now().UTC()
	req := &domain.EditAccessRequest{
		ID:        uuid.New(),
		UserID:    userID,
		Message:   strings.TrimSpace(message),
		Status:    domain.AccessRequestPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.requests.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *accessRequestService) MyLatest(ctx context.Context, userID uuid.UUID) (*domain.EditAccessRequest, error) {
	return s.requests.GetLatestForUser(ctx, userID)
}

func (s *accessRequestService) ListPending(ctx context.Context, adminID uuid.UUID) ([]domain.AccessRequestSummary, error) {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	reqs, err := s.requests.ListByStatus(ctx, domain.AccessRequestPending)
	if err != nil {
		return nil, err
	}
	out := make([]domain.AccessRequestSummary, 0, len(reqs))
	for _, req := range reqs {
		user, err := s.users.GetByID(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		out = append(out, domain.AccessRequestSummary{Request: req, User: *user})
	}
	return out, nil
}

func (s *accessRequestService) Approve(ctx context.Context, adminID, requestID uuid.UUID) (*domain.EditAccessRequest, error) {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req.Status != domain.AccessRequestPending {
		return nil, domain.ErrForbidden
	}

	familyID, err := s.royal.ID(ctx)
	if err != nil {
		return nil, err
	}

	_, err = s.collaborators.Get(ctx, familyID, req.UserID)
	if errors.Is(err, domain.ErrNotFound) {
		collab := &domain.FamilyCollaborator{
			ID:       uuid.New(),
			FamilyID: familyID,
			UserID:   req.UserID,
			Role:     domain.FamilyRoleEditor,
		}
		if err := s.collaborators.Create(ctx, collab); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	req.Status = domain.AccessRequestApproved
	req.ReviewedBy = &adminID
	req.ReviewedAt = &now
	req.UpdatedAt = now
	if err := s.requests.Update(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *accessRequestService) Reject(ctx context.Context, adminID, requestID uuid.UUID) (*domain.EditAccessRequest, error) {
	if err := s.requireAdmin(ctx, adminID); err != nil {
		return nil, err
	}
	req, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if req.Status != domain.AccessRequestPending {
		return nil, domain.ErrForbidden
	}
	now := time.Now().UTC()
	req.Status = domain.AccessRequestRejected
	req.ReviewedBy = &adminID
	req.ReviewedAt = &now
	req.UpdatedAt = now
	if err := s.requests.Update(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *accessRequestService) requireAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsAdmin() {
		return domain.ErrForbidden
	}
	return nil
}

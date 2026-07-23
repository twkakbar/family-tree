package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

func TestAccessRequestApproveGrantsEditor(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	requests := newFakeAccessRequestRepo()

	adminID := uuid.New()
	userID := uuid.New()
	_ = users.Create(context.Background(), &domain.User{ID: adminID, Name: "Admin", Email: "a@ex.com", Role: domain.RoleAdmin})
	_ = users.Create(context.Background(), &domain.User{ID: userID, Name: "Desc", Email: "d@ex.com", Role: domain.RoleUser})

	royal := NewRoyalFamilyService(families, collabs, users, "Kesultanan Aceh")
	family, err := royal.Ensure(context.Background(), adminID, "Kesultanan Aceh", "desc")
	if err != nil {
		t.Fatal(err)
	}
	collabs.seedOwner(family.ID, adminID)

	svc := NewAccessRequestService(requests, users, royal, collabs)
	req, err := svc.Submit(context.Background(), userID, "I am a living descendant of the sultan line.")
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}

	approved, err := svc.Approve(context.Background(), adminID, req.ID)
	if err != nil {
		t.Fatalf("Approve() error = %v", err)
	}
	if approved.Status != domain.AccessRequestApproved {
		t.Errorf("status = %q", approved.Status)
	}

	can, err := royal.CanEdit(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if !can {
		t.Error("expected editor grant after approve")
	}
}

func TestAccessRequestRejectDoesNotGrant(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	requests := newFakeAccessRequestRepo()

	adminID := uuid.New()
	userID := uuid.New()
	_ = users.Create(context.Background(), &domain.User{ID: adminID, Name: "Admin", Email: "a@ex.com", Role: domain.RoleAdmin})
	_ = users.Create(context.Background(), &domain.User{ID: userID, Name: "Desc", Email: "d@ex.com", Role: domain.RoleUser})

	royal := NewRoyalFamilyService(families, collabs, users, "Kesultanan Aceh")
	family, _ := royal.Ensure(context.Background(), adminID, "Kesultanan Aceh", "desc")
	collabs.seedOwner(family.ID, adminID)

	svc := NewAccessRequestService(requests, users, royal, collabs)
	req, _ := svc.Submit(context.Background(), userID, "Please let me edit the royal tree.")
	_, err := svc.Reject(context.Background(), adminID, req.ID)
	if err != nil {
		t.Fatal(err)
	}
	can, _ := royal.CanEdit(context.Background(), userID)
	if can {
		t.Error("reject must not grant edit access")
	}
}

func TestDuplicatePendingRequest(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	requests := newFakeAccessRequestRepo()

	adminID := uuid.New()
	userID := uuid.New()
	_ = users.Create(context.Background(), &domain.User{ID: adminID, Name: "Admin", Email: "a@ex.com", Role: domain.RoleAdmin})
	_ = users.Create(context.Background(), &domain.User{ID: userID, Name: "Desc", Email: "d@ex.com", Role: domain.RoleUser})

	royal := NewRoyalFamilyService(families, collabs, users, "Kesultanan Aceh")
	family, _ := royal.Ensure(context.Background(), adminID, "Kesultanan Aceh", "desc")
	collabs.seedOwner(family.ID, adminID)

	svc := NewAccessRequestService(requests, users, royal, collabs)
	_, _ = svc.Submit(context.Background(), userID, "First request with enough detail.")
	_, err := svc.Submit(context.Background(), userID, "Second request with enough detail.")
	if !errors.Is(err, domain.ErrAccessRequestPending) {
		t.Errorf("error = %v, want ErrAccessRequestPending", err)
	}
}

type fakeAccessRequestRepo struct {
	items map[uuid.UUID]*domain.EditAccessRequest
}

func newFakeAccessRequestRepo() *fakeAccessRequestRepo {
	return &fakeAccessRequestRepo{items: make(map[uuid.UUID]*domain.EditAccessRequest)}
}

func (f *fakeAccessRequestRepo) Create(_ context.Context, req *domain.EditAccessRequest) error {
	for _, existing := range f.items {
		if existing.UserID == req.UserID && existing.Status == domain.AccessRequestPending {
			return domain.ErrAccessRequestPending
		}
	}
	cp := *req
	f.items[req.ID] = &cp
	return nil
}

func (f *fakeAccessRequestRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.EditAccessRequest, error) {
	if req, ok := f.items[id]; ok {
		cp := *req
		return &cp, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeAccessRequestRepo) GetLatestForUser(_ context.Context, userID uuid.UUID) (*domain.EditAccessRequest, error) {
	var latest *domain.EditAccessRequest
	for _, req := range f.items {
		if req.UserID != userID {
			continue
		}
		if latest == nil || req.CreatedAt.After(latest.CreatedAt) {
			cp := *req
			latest = &cp
		}
	}
	if latest == nil {
		return nil, domain.ErrNotFound
	}
	return latest, nil
}

func (f *fakeAccessRequestRepo) ListByStatus(_ context.Context, status domain.AccessRequestStatus) ([]domain.EditAccessRequest, error) {
	out := make([]domain.EditAccessRequest, 0)
	for _, req := range f.items {
		if req.Status == status {
			out = append(out, *req)
		}
	}
	return out, nil
}

func (f *fakeAccessRequestRepo) Update(_ context.Context, req *domain.EditAccessRequest) error {
	if _, ok := f.items[req.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *req
	f.items[req.ID] = &cp
	return nil
}

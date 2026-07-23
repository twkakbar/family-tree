package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

func TestCreateFamily(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	svc := NewFamilyService(families, collabs, users)

	ownerID := uuid.New()
	ctx := context.Background()

	summary, err := svc.Create(ctx, ownerID, domain.CreateFamilyInput{
		Name:        "  The Does  ",
		Description: " Paternal branch ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if summary.Family.Name != "The Does" {
		t.Errorf("name = %q, want trimmed", summary.Family.Name)
	}
	if summary.Role != domain.FamilyRoleOwner {
		t.Errorf("role = %q, want owner", summary.Role)
	}
}

func TestAddCollaborator(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	svc := NewFamilyService(families, collabs, users)
	ctx := context.Background()

	ownerID := uuid.New()
	invitee := &domain.User{ID: uuid.New(), Name: "Jane", Email: "jane@example.com"}
	if err := users.Create(ctx, invitee); err != nil {
		t.Fatal(err)
	}

	created, err := svc.Create(ctx, ownerID, domain.CreateFamilyInput{Name: "Family"})
	if err != nil {
		t.Fatal(err)
	}
	collabs.seedOwner(created.Family.ID, ownerID)

	added, err := svc.AddCollaborator(ctx, ownerID, created.Family.ID, domain.AddCollaboratorInput{
		Email: "jane@example.com",
		Role:  domain.FamilyRoleEditor,
	})
	if err != nil {
		t.Fatalf("AddCollaborator() error = %v", err)
	}
	if added.Collaborator.Role != domain.FamilyRoleEditor {
		t.Errorf("role = %q, want editor", added.Collaborator.Role)
	}

	families.grant(invitee.ID, domain.FamilySummary{
		Family:       created.Family,
		Role:         domain.FamilyRoleEditor,
		PersonsCount: 0,
	})

	_, err = svc.Update(ctx, invitee.ID, created.Family.ID, domain.UpdateFamilyInput{Name: "Renamed"})
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("editor update error = %v, want ErrForbidden", err)
	}
}

func TestRemoveLastOwner(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	svc := NewFamilyService(families, collabs, users)
	ctx := context.Background()

	ownerID := uuid.New()
	created, err := svc.Create(ctx, ownerID, domain.CreateFamilyInput{Name: "Family"})
	if err != nil {
		t.Fatal(err)
	}
	collabs.seedOwner(created.Family.ID, ownerID)

	err = svc.RemoveCollaborator(ctx, ownerID, created.Family.ID, ownerID)
	if !errors.Is(err, domain.ErrLastOwner) {
		t.Errorf("RemoveCollaborator() error = %v, want ErrLastOwner", err)
	}
}

func TestViewerCannotDeleteFamily(t *testing.T) {
	users := newFakeUserRepo()
	families := newFakeFamilyRepo()
	collabs := newFakeFamilyCollaboratorRepo()
	svc := NewFamilyService(families, collabs, users)
	ctx := context.Background()

	ownerID := uuid.New()
	viewerID := uuid.New()
	familyID := uuid.New()
	now := time.Now().UTC()
	family := domain.Family{ID: familyID, Name: "Family", CreatedBy: ownerID, CreatedAt: now, UpdatedAt: now}
	families.families[familyID] = &family
	families.grant(ownerID, domain.FamilySummary{Family: family, Role: domain.FamilyRoleOwner})
	families.grant(viewerID, domain.FamilySummary{Family: family, Role: domain.FamilyRoleViewer})

	err := svc.Delete(ctx, viewerID, familyID)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("Delete() error = %v, want ErrForbidden", err)
	}
}

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

func seedEditor(t *testing.T, users *fakeUserRepo, families *fakeFamilyRepo, role domain.FamilyRole) (uuid.UUID, uuid.UUID) {
	t.Helper()
	userID := uuid.New()
	familyID := uuid.New()
	now := time.Now().UTC()
	if err := users.Create(context.Background(), &domain.User{
		ID: userID, Name: "Editor", Email: userID.String() + "@example.com", Role: domain.RoleUser,
	}); err != nil {
		t.Fatal(err)
	}
	families.families[familyID] = &domain.Family{ID: familyID, Name: "Family", CreatedAt: now, UpdatedAt: now}
	families.grant(userID, domain.FamilySummary{
		Family: *families.families[familyID],
		Role:   role,
	})
	return userID, familyID
}

func TestCreatePerson(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	svc := NewPersonService(families, persons, rels, users, t.TempDir(), 1<<20)

	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)

	person, err := svc.Create(context.Background(), userID, familyID, domain.CreatePersonInput{
		FirstName: "John",
		LastName:  "Doe",
		Gender:    domain.GenderMale,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if person.FirstName != "John" {
		t.Errorf("first_name = %q", person.FirstName)
	}
}

func TestCreateParentRelationship(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	svc := NewPersonService(families, persons, rels, users, t.TempDir(), 1<<20)
	ctx := context.Background()

	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)

	parent, err := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{
		FirstName: "Parent",
		Gender:    domain.GenderMale,
	})
	if err != nil {
		t.Fatal(err)
	}
	child, err := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{
		FirstName: "Child",
		Gender:    domain.GenderFemale,
	})
	if err != nil {
		t.Fatal(err)
	}

	rel, err := svc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: parent.ID,
		ToPersonID:   child.ID,
		Type:         domain.RelationshipParent,
	})
	if err != nil {
		t.Fatalf("CreateRelationship() error = %v", err)
	}
	if rel.Type != domain.RelationshipParent {
		t.Errorf("type = %q", rel.Type)
	}
}

func TestRejectParentCycle(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	svc := NewPersonService(families, persons, rels, users, t.TempDir(), 1<<20)
	ctx := context.Background()

	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)

	a, _ := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{FirstName: "A", Gender: domain.GenderMale})
	b, _ := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{FirstName: "B", Gender: domain.GenderMale})
	_, _ = svc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: a.ID, ToPersonID: b.ID, Type: domain.RelationshipParent,
	})

	_, err := svc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: b.ID, ToPersonID: a.ID, Type: domain.RelationshipParent,
	})
	if !errors.Is(err, domain.ErrInvalidRelationship) {
		t.Errorf("error = %v, want ErrInvalidRelationship", err)
	}
}

func TestRejectThirdParent(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	svc := NewPersonService(families, persons, rels, users, t.TempDir(), 1<<20)
	ctx := context.Background()

	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)

	child, _ := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{FirstName: "Child", Gender: domain.GenderUnknown})
	p1, _ := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{FirstName: "P1", Gender: domain.GenderMale})
	p2, _ := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{FirstName: "P2", Gender: domain.GenderFemale})
	p3, _ := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{FirstName: "P3", Gender: domain.GenderMale})

	_, _ = svc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: p1.ID, ToPersonID: child.ID, Type: domain.RelationshipParent,
	})
	_, _ = svc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: p2.ID, ToPersonID: child.ID, Type: domain.RelationshipParent,
	})
	_, err := svc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: p3.ID, ToPersonID: child.ID, Type: domain.RelationshipParent,
	})
	if !errors.Is(err, domain.ErrInvalidRelationship) {
		t.Errorf("error = %v, want ErrInvalidRelationship", err)
	}
}

func TestUpdatePhoto(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	dir := t.TempDir()
	svc := NewPersonService(families, persons, rels, users, dir, 1<<20)
	ctx := context.Background()

	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)
	person, err := svc.Create(ctx, userID, familyID, domain.CreatePersonInput{
		FirstName: "Photo", Gender: domain.GenderFemale,
	})
	if err != nil {
		t.Fatal(err)
	}

	pngHeader := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d}
	updated, err := svc.UpdatePhoto(ctx, userID, person.ID, pngHeader)
	if err != nil {
		t.Fatalf("UpdatePhoto() error = %v", err)
	}
	if updated.PhotoPath == "" {
		t.Fatal("expected photo_path set")
	}

	_, err = svc.UpdatePhoto(ctx, userID, person.ID, []byte("not-an-image"))
	if !errors.Is(err, domain.ErrInvalidUpload) {
		t.Errorf("error = %v, want ErrInvalidUpload", err)
	}
}

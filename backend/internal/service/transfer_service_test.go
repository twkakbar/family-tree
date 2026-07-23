package service

import (
	"context"
	"errors"
	"testing"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

func TestTransferExportImportRoundTrip(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	personSvc := NewPersonService(families, persons, rels, users, t.TempDir(), 1<<20)
	familySvc := NewFamilyService(families, newFakeFamilyCollaboratorRepo(), users)
	transfer := NewTransferService(familySvc, personSvc)
	ctx := context.Background()

	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)

	parent, err := personSvc.Create(ctx, userID, familyID, domain.CreatePersonInput{
		FirstName: "Parent", Gender: domain.GenderMale,
	})
	if err != nil {
		t.Fatal(err)
	}
	child, err := personSvc.Create(ctx, userID, familyID, domain.CreatePersonInput{
		FirstName: "Child", Gender: domain.GenderFemale,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = personSvc.CreateRelationship(ctx, userID, familyID, domain.CreateRelationshipInput{
		FromPersonID: parent.ID, ToPersonID: child.ID, Type: domain.RelationshipParent,
	})
	if err != nil {
		t.Fatal(err)
	}

	doc, err := transfer.ExportJSON(ctx, userID, familyID)
	if err != nil {
		t.Fatalf("ExportJSON: %v", err)
	}
	if doc.FormatVersion != domain.TreeTransferFormatVersion {
		t.Fatalf("version = %d", doc.FormatVersion)
	}
	if len(doc.Persons) != 2 || len(doc.Relationships) != 1 {
		t.Fatalf("export size persons=%d rels=%d", len(doc.Persons), len(doc.Relationships))
	}

	result, err := transfer.ImportJSON(ctx, userID, familyID, doc)
	if err != nil {
		t.Fatalf("ImportJSON: %v", err)
	}
	if result.PersonsCreated != 2 || result.RelationshipsCreated != 1 {
		t.Fatalf("import result = %+v", result)
	}

	tree, err := personSvc.GetFamilyTree(ctx, userID, familyID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tree.Persons) != 4 || len(tree.Relationships) != 2 {
		t.Fatalf("after import persons=%d rels=%d", len(tree.Persons), len(tree.Relationships))
	}
}

func TestTransferRejectsBadVersion(t *testing.T) {
	families := newFakeFamilyRepo()
	persons := newFakePersonRepo()
	rels := newFakeRelationshipRepo()
	users := newFakeUserRepo()
	personSvc := NewPersonService(families, persons, rels, users, t.TempDir(), 1<<20)
	familySvc := NewFamilyService(families, newFakeFamilyCollaboratorRepo(), users)
	transfer := NewTransferService(familySvc, personSvc)
	userID, familyID := seedEditor(t, users, families, domain.FamilyRoleEditor)

	_, err := transfer.ImportJSON(context.Background(), userID, familyID, &domain.TreeTransfer{
		FormatVersion: 99,
	})
	if !errors.Is(err, domain.ErrUnsupportedFormat) {
		t.Fatalf("error = %v", err)
	}
}

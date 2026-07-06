package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

func TestUpdateProfile(t *testing.T) {
	users := newFakeUserRepo()
	svc := NewUserService(users)
	ctx := context.Background()

	user := &domain.User{ID: uuid.New(), Name: "Old Name", Email: "u@example.com"}
	if err := users.Create(ctx, user); err != nil {
		t.Fatal(err)
	}

	updated, err := svc.UpdateProfile(ctx, user.ID, domain.UpdateProfileInput{Name: "  New Name  "})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("name = %q, want trimmed %q", updated.Name, "New Name")
	}
}

func TestUpdateProfileNotFound(t *testing.T) {
	svc := NewUserService(newFakeUserRepo())
	_, err := svc.UpdateProfile(context.Background(), uuid.New(), domain.UpdateProfileInput{Name: "X"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

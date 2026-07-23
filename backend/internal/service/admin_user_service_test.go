package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/password"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

func TestPasswordResetFlow(t *testing.T) {
	users := newFakeUserRepo()
	tokens := newFakeRefreshTokenRepo()
	resets := newFakePasswordResetTokenRepo()
	jwt := token.NewJWTManager("test-secret-that-is-32-characters!!", "test", 15*time.Minute)
	hasher := password.NewHasher(4)
	svc := NewAuthService(users, tokens, resets, jwt, hasher, time.Hour)

	_, err := svc.Register(context.Background(), domain.RegisterInput{
		Name: "Reset User", Email: "reset@example.com", Password: "password123",
	})
	if err != nil {
		t.Fatal(err)
	}

	raw, err := svc.RequestPasswordReset(context.Background(), "reset@example.com")
	if err != nil || raw == "" {
		t.Fatalf("RequestPasswordReset raw=%q err=%v", raw, err)
	}

	if err := svc.ResetPassword(context.Background(), raw, "newpass99"); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}

	_, err = svc.Login(context.Background(), domain.LoginInput{
		Email: "reset@example.com", Password: "newpass99",
	})
	if err != nil {
		t.Fatalf("login after reset: %v", err)
	}

	err = svc.ResetPassword(context.Background(), raw, "another99")
	if !errors.Is(err, domain.ErrTokenInvalid) && !errors.Is(err, domain.ErrTokenExpired) {
		t.Fatalf("reuse token error = %v", err)
	}
}

func TestAdminListAndDeleteUser(t *testing.T) {
	users := newFakeUserRepo()
	tokens := newFakeRefreshTokenRepo()
	resets := newFakePasswordResetTokenRepo()
	adminSvc := NewAdminUserService(users, tokens, resets)

	adminID := uuid.New()
	userID := uuid.New()
	_ = users.Create(context.Background(), &domain.User{
		ID: adminID, Name: "Admin", Email: "admin@example.com", Role: domain.RoleAdmin,
	})
	_ = users.Create(context.Background(), &domain.User{
		ID: userID, Name: "Member", Email: "member@example.com", Role: domain.RoleUser,
	})

	list, total, err := adminSvc.List(context.Background(), adminID, "", 1, 20)
	if err != nil || total < 2 || len(list) < 2 {
		t.Fatalf("List total=%d len=%d err=%v", total, len(list), err)
	}

	if err := adminSvc.Delete(context.Background(), adminID, userID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := adminSvc.Delete(context.Background(), adminID, adminID); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("self-delete error = %v", err)
	}
}

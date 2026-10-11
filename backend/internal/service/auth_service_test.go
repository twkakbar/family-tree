package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/password"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

func newTestAuthService() (domain.AuthService, *fakeUserRepo, *fakeRefreshTokenRepo) {
	users := newFakeUserRepo()
	tokens := newFakeRefreshTokenRepo()
	resets := newFakePasswordResetTokenRepo()
	jwt := token.NewJWTManager("test-secret-that-is-32-characters!!", "test", 15*time.Minute)
	hasher := password.NewHasher(4) // low cost for fast tests
	svc := NewAuthService(users, tokens, resets, nil, nil, jwt, hasher, time.Hour, nil, "")
	return svc, users, tokens
}

func register(t *testing.T, svc domain.AuthService) *domain.AuthResult {
	t.Helper()
	result, err := svc.Register(context.Background(), domain.RegisterInput{
		Name:     "Jane Doe",
		Email:    "Jane@Example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if result.Session == nil {
		t.Fatal("Register() session is nil")
	}
	return result.Session
}

func TestRegister(t *testing.T) {
	svc, _, _ := newTestAuthService()
	result := register(t, svc)

	if result.User.Email != "jane@example.com" {
		t.Errorf("email not normalized: %q", result.User.Email)
	}
	if result.User.Role != domain.RoleUser {
		t.Errorf("role = %q, want %q", result.User.Role, domain.RoleUser)
	}
	if result.User.PasswordHash == "password123" || result.User.PasswordHash == "" {
		t.Error("password must be stored hashed")
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Error("registration must establish a session")
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	svc, _, _ := newTestAuthService()
	register(t, svc)

	_, err := svc.Register(context.Background(), domain.RegisterInput{
		Name:     "Other",
		Email:    "JANE@example.com",
		Password: "password456",
	})
	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Errorf("error = %v, want ErrEmailTaken", err)
	}
}

func TestLogin(t *testing.T) {
	svc, _, _ := newTestAuthService()
	register(t, svc)

	result, err := svc.Login(context.Background(), domain.LoginInput{
		Email:    "jane@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Error("login must return a token pair")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	svc, _, _ := newTestAuthService()
	register(t, svc)

	_, err := svc.Login(context.Background(), domain.LoginInput{
		Email:    "jane@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("error = %v, want ErrInvalidCredentials", err)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	svc, _, _ := newTestAuthService()

	_, err := svc.Login(context.Background(), domain.LoginInput{
		Email:    "nobody@example.com",
		Password: "password123",
	})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("error = %v, want ErrInvalidCredentials (must not leak account existence)", err)
	}
}

func TestRefreshRotatesToken(t *testing.T) {
	svc, _, tokens := newTestAuthService()
	session := register(t, svc)

	rotated, err := svc.Refresh(context.Background(), session.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if rotated.RefreshToken == session.RefreshToken {
		t.Error("refresh token must rotate")
	}
	if got := tokens.activeCountForUser(session.User.ID); got != 1 {
		t.Errorf("active tokens after rotation = %d, want 1", got)
	}
}

func TestRefreshReuseDetection(t *testing.T) {
	svc, _, tokens := newTestAuthService()
	session := register(t, svc)

	// First refresh consumes the token…
	if _, err := svc.Refresh(context.Background(), session.RefreshToken); err != nil {
		t.Fatalf("first Refresh() error = %v", err)
	}
	// …presenting it again is treated as theft: all sessions are revoked.
	_, err := svc.Refresh(context.Background(), session.RefreshToken)
	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Fatalf("error = %v, want ErrTokenInvalid", err)
	}
	if got := tokens.activeCountForUser(session.User.ID); got != 0 {
		t.Errorf("active tokens after reuse detection = %d, want 0", got)
	}
}

func TestRefreshUnknownToken(t *testing.T) {
	svc, _, _ := newTestAuthService()
	_, err := svc.Refresh(context.Background(), "unknown-token")
	if !errors.Is(err, domain.ErrTokenInvalid) {
		t.Errorf("error = %v, want ErrTokenInvalid", err)
	}
}

func TestLogoutIsIdempotent(t *testing.T) {
	svc, _, tokens := newTestAuthService()
	session := register(t, svc)

	for i := 0; i < 2; i++ {
		if err := svc.Logout(context.Background(), session.RefreshToken); err != nil {
			t.Fatalf("Logout() call %d error = %v", i+1, err)
		}
	}
	if err := svc.Logout(context.Background(), ""); err != nil {
		t.Fatalf("Logout(empty) error = %v", err)
	}
	if got := tokens.activeCountForUser(session.User.ID); got != 0 {
		t.Errorf("active tokens after logout = %d, want 0", got)
	}
}

func TestEnsureAdmin(t *testing.T) {
	svc, users, _ := newTestAuthService()
	ctx := context.Background()

	if err := svc.EnsureAdmin(ctx, "Admin", "admin@example.com", "secret-admin-pw1"); err != nil {
		t.Fatalf("EnsureAdmin() error = %v", err)
	}
	admin, err := users.GetByEmail(ctx, "admin@example.com")
	if err != nil {
		t.Fatalf("admin not created: %v", err)
	}
	if admin.Role != domain.RoleAdmin {
		t.Errorf("role = %q, want admin", admin.Role)
	}

	// Second call must not modify the existing account.
	if err := svc.EnsureAdmin(ctx, "Admin", "admin@example.com", "different-pw"); err != nil {
		t.Fatalf("EnsureAdmin() second call error = %v", err)
	}
	again, _ := users.GetByEmail(ctx, "admin@example.com")
	if again.PasswordHash != admin.PasswordHash {
		t.Error("EnsureAdmin must never overwrite an existing password")
	}
}

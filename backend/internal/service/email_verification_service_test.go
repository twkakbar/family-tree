package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
	"github.com/twkakbar/family-tree/backend/pkg/password"
	"github.com/twkakbar/family-tree/backend/pkg/token"
)

type recordingMailer struct {
	sends    int
	lastText string
}

func (m *recordingMailer) Enabled() bool { return true }
func (m *recordingMailer) Send(_ context.Context, _, _, text string) error {
	m.sends++
	m.lastText = text
	return nil
}

type memVerifyRepo struct {
	byHash map[string]*domain.EmailVerificationToken
}

func (r *memVerifyRepo) Create(_ context.Context, item *domain.EmailVerificationToken) error {
	if r.byHash == nil {
		r.byHash = map[string]*domain.EmailVerificationToken{}
	}
	cp := *item
	r.byHash[item.TokenHash] = &cp
	return nil
}

func (r *memVerifyRepo) GetByHash(_ context.Context, hash string) (*domain.EmailVerificationToken, error) {
	item, ok := r.byHash[hash]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *item
	return &cp, nil
}

func (r *memVerifyRepo) Revoke(_ context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	for _, item := range r.byHash {
		if item.ID == id {
			item.RevokedAt = &now
		}
	}
	return nil
}

func (r *memVerifyRepo) RevokeAllForUser(_ context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	for _, item := range r.byHash {
		if item.UserID == userID && item.RevokedAt == nil {
			item.RevokedAt = &now
		}
	}
	return nil
}

type memDispatchRepo struct {
	rows []domain.EmailDispatch
}

func (r *memDispatchRepo) Create(_ context.Context, row *domain.EmailDispatch) error {
	r.rows = append(r.rows, *row)
	return nil
}

func (r *memDispatchRepo) CountSince(_ context.Context, email, purpose string, since time.Time) (int64, error) {
	var n int64
	for _, row := range r.rows {
		if row.Email == email && row.Purpose == purpose && !row.SentAt.Before(since) {
			n++
		}
	}
	return n, nil
}

func TestRegisterSendsVerificationAndRateLimits(t *testing.T) {
	users := newFakeUserRepo()
	tokens := newFakeRefreshTokenRepo()
	resets := newFakePasswordResetTokenRepo()
	verifies := &memVerifyRepo{}
	dispatches := &memDispatchRepo{}
	mailer := &recordingMailer{}
	jwt := token.NewJWTManager("test-secret-that-is-32-characters!!", "test", 15*time.Minute)
	hasher := password.NewHasher(4)
	svc := NewAuthService(users, tokens, resets, verifies, dispatches, jwt, hasher, time.Hour, mailer, "https://twkakbar.my.id")

	for i := 0; i < 3; i++ {
		out, err := svc.Register(context.Background(), domain.RegisterInput{
			Name: "Verify User", Email: "verify@example.com", Password: "password123",
		})
		if err != nil {
			t.Fatalf("send %d: %v", i+1, err)
		}
		if !out.VerificationRequired || out.Session != nil {
			t.Fatalf("send %d: want verification without session", i+1)
		}
	}
	if mailer.sends != 3 {
		t.Fatalf("sends = %d, want 3", mailer.sends)
	}

	_, err := svc.Register(context.Background(), domain.RegisterInput{
		Name: "Verify User", Email: "verify@example.com", Password: "password123",
	})
	if !errors.Is(err, domain.ErrRateLimited) {
		t.Fatalf("fourth send error = %v, want rate limited", err)
	}

	_, err = svc.Login(context.Background(), domain.LoginInput{
		Email: "verify@example.com", Password: "password123",
	})
	if !errors.Is(err, domain.ErrEmailNotVerified) {
		t.Fatalf("login error = %v, want unverified", err)
	}

	raw := verificationToken(t, mailer.lastText)
	if err := svc.VerifyEmail(context.Background(), raw); err != nil {
		t.Fatalf("VerifyEmail: %v", err)
	}
	if _, err := svc.Login(context.Background(), domain.LoginInput{
		Email: "verify@example.com", Password: "password123",
	}); err != nil {
		t.Fatalf("login after verify: %v", err)
	}
}

func verificationToken(t *testing.T, body string) string {
	t.Helper()
	const marker = "token="
	i := strings.LastIndex(body, marker)
	if i < 0 {
		t.Fatalf("verification link missing from %q", body)
	}
	token := body[i+len(marker):]
	if nl := strings.IndexAny(token, "\r\n "); nl >= 0 {
		token = token[:nl]
	}
	decoded, err := url.QueryUnescape(token)
	if err != nil {
		t.Fatalf("unescape token: %v", err)
	}
	return decoded
}

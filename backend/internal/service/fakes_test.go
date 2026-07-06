package service

// In-memory fakes of the domain repository interfaces. This is the reference
// pattern for service unit tests: no mocking framework, no database.

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/twkakbar/family-tree/backend/internal/domain"
)

type fakeUserRepo struct {
	mu    sync.Mutex
	users map[uuid.UUID]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[uuid.UUID]*domain.User)}
}

func (f *fakeUserRepo) Create(_ context.Context, user *domain.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.users {
		if u.Email == user.Email {
			return domain.ErrEmailTaken
		}
	}
	cp := *user
	f.users[user.ID] = &cp
	return nil
}

func (f *fakeUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if u, ok := f.users[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, u := range f.users {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepo) Update(_ context.Context, user *domain.User) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.users[user.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *user
	f.users[user.ID] = &cp
	return nil
}

type fakeRefreshTokenRepo struct {
	mu     sync.Mutex
	tokens map[uuid.UUID]*domain.RefreshToken
}

func newFakeRefreshTokenRepo() *fakeRefreshTokenRepo {
	return &fakeRefreshTokenRepo{tokens: make(map[uuid.UUID]*domain.RefreshToken)}
}

func (f *fakeRefreshTokenRepo) Create(_ context.Context, token *domain.RefreshToken) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *token
	f.tokens[token.ID] = &cp
	return nil
}

func (f *fakeRefreshTokenRepo) GetByHash(_ context.Context, hash string) (*domain.RefreshToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, t := range f.tokens {
		if t.TokenHash == hash {
			cp := *t
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeRefreshTokenRepo) Revoke(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t, ok := f.tokens[id]; ok && t.RevokedAt == nil {
		now := time.Now().UTC()
		t.RevokedAt = &now
	}
	return nil
}

func (f *fakeRefreshTokenRepo) RevokeAllForUser(_ context.Context, userID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now().UTC()
	for _, t := range f.tokens {
		if t.UserID == userID && t.RevokedAt == nil {
			t.RevokedAt = &now
		}
	}
	return nil
}

func (f *fakeRefreshTokenRepo) activeCountForUser(userID uuid.UUID) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, t := range f.tokens {
		if t.UserID == userID && t.RevokedAt == nil {
			n++
		}
	}
	return n
}

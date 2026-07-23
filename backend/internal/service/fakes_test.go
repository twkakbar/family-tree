package service

// In-memory fakes of the domain repository interfaces. This is the reference
// pattern for service unit tests: no mocking framework, no database.

import (
	"context"
	"strings"
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

func (f *fakeUserRepo) List(_ context.Context, search string, offset, limit int) ([]domain.User, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var all []domain.User
	for _, u := range f.users {
		if search != "" {
			if !containsFold(u.Name, search) && !containsFold(u.Email, search) {
				continue
			}
		}
		all = append(all, *u)
	}
	total := int64(len(all))
	if offset > len(all) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], total, nil
}

func (f *fakeUserRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.users[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.users, id)
	return nil
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
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

type fakePasswordResetTokenRepo struct {
	mu     sync.Mutex
	tokens map[uuid.UUID]*domain.PasswordResetToken
}

func newFakePasswordResetTokenRepo() *fakePasswordResetTokenRepo {
	return &fakePasswordResetTokenRepo{tokens: make(map[uuid.UUID]*domain.PasswordResetToken)}
}

func (f *fakePasswordResetTokenRepo) Create(_ context.Context, token *domain.PasswordResetToken) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *token
	f.tokens[token.ID] = &cp
	return nil
}

func (f *fakePasswordResetTokenRepo) GetByHash(_ context.Context, hash string) (*domain.PasswordResetToken, error) {
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

func (f *fakePasswordResetTokenRepo) Revoke(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if t, ok := f.tokens[id]; ok && t.RevokedAt == nil {
		now := time.Now().UTC()
		t.RevokedAt = &now
	}
	return nil
}

func (f *fakePasswordResetTokenRepo) RevokeAllForUser(_ context.Context, userID uuid.UUID) error {
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

type fakeFamilyRepo struct {
	mu        sync.Mutex
	families  map[uuid.UUID]*domain.Family
	summaries map[uuid.UUID]map[uuid.UUID]*domain.FamilySummary // userID -> familyID
}

func newFakeFamilyRepo() *fakeFamilyRepo {
	return &fakeFamilyRepo{
		families:  make(map[uuid.UUID]*domain.Family),
		summaries: make(map[uuid.UUID]map[uuid.UUID]*domain.FamilySummary),
	}
}

func (f *fakeFamilyRepo) Create(_ context.Context, family *domain.Family) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *family
	f.families[family.ID] = &cp
	return nil
}

func (f *fakeFamilyRepo) CreateWithOwner(ctx context.Context, family *domain.Family, ownerID uuid.UUID, role domain.FamilyRole) error {
	if err := f.Create(ctx, family); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.summaries[ownerID] == nil {
		f.summaries[ownerID] = make(map[uuid.UUID]*domain.FamilySummary)
	}
	f.summaries[ownerID][family.ID] = &domain.FamilySummary{
		Family:       *f.families[family.ID],
		Role:         role,
		PersonsCount: 0,
	}
	return nil
}

func (f *fakeFamilyRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Family, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if family, ok := f.families[id]; ok {
		cp := *family
		return &cp, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeFamilyRepo) GetByName(_ context.Context, name string) (*domain.Family, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, family := range f.families {
		if family.Name == name {
			cp := *family
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeFamilyRepo) Update(_ context.Context, family *domain.Family) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.families[family.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *family
	f.families[family.ID] = &cp
	for _, byFamily := range f.summaries {
		if summary, ok := byFamily[family.ID]; ok {
			summary.Family = cp
		}
	}
	return nil
}

func (f *fakeFamilyRepo) SoftDelete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.families[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.families, id)
	for userID := range f.summaries {
		delete(f.summaries[userID], id)
	}
	return nil
}

func (f *fakeFamilyRepo) ListForUser(_ context.Context, userID uuid.UUID, offset, limit int) ([]domain.FamilySummary, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	byFamily := f.summaries[userID]
	items := make([]domain.FamilySummary, 0, len(byFamily))
	for _, summary := range byFamily {
		items = append(items, *summary)
	}
	total := int64(len(items))
	if offset >= len(items) {
		return []domain.FamilySummary{}, total, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], total, nil
}

func (f *fakeFamilyRepo) GetSummaryForUser(_ context.Context, userID, familyID uuid.UUID) (*domain.FamilySummary, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if byFamily, ok := f.summaries[userID]; ok {
		if summary, ok := byFamily[familyID]; ok {
			cp := *summary
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeFamilyRepo) grant(userID uuid.UUID, summary domain.FamilySummary) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.summaries[userID] == nil {
		f.summaries[userID] = make(map[uuid.UUID]*domain.FamilySummary)
	}
	cp := summary
	f.summaries[userID][summary.Family.ID] = &cp
}

type fakeFamilyCollaboratorRepo struct {
	mu      sync.Mutex
	collabs map[uuid.UUID]map[uuid.UUID]*domain.FamilyCollaborator // familyID -> userID
}

func newFakeFamilyCollaboratorRepo() *fakeFamilyCollaboratorRepo {
	return &fakeFamilyCollaboratorRepo{collabs: make(map[uuid.UUID]map[uuid.UUID]*domain.FamilyCollaborator)}
}

func (f *fakeFamilyCollaboratorRepo) Create(_ context.Context, collab *domain.FamilyCollaborator) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.collabs[collab.FamilyID] == nil {
		f.collabs[collab.FamilyID] = make(map[uuid.UUID]*domain.FamilyCollaborator)
	}
	if _, exists := f.collabs[collab.FamilyID][collab.UserID]; exists {
		return domain.ErrAlreadyCollaborator
	}
	cp := *collab
	f.collabs[collab.FamilyID][collab.UserID] = &cp
	return nil
}

func (f *fakeFamilyCollaboratorRepo) Get(_ context.Context, familyID, userID uuid.UUID) (*domain.FamilyCollaborator, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if byUser, ok := f.collabs[familyID]; ok {
		if collab, ok := byUser[userID]; ok {
			cp := *collab
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (f *fakeFamilyCollaboratorRepo) ListByFamily(_ context.Context, familyID uuid.UUID) ([]domain.FamilyCollaborator, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]domain.FamilyCollaborator, 0)
	for _, collab := range f.collabs[familyID] {
		out = append(out, *collab)
	}
	return out, nil
}

func (f *fakeFamilyCollaboratorRepo) Update(_ context.Context, collab *domain.FamilyCollaborator) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if byUser, ok := f.collabs[collab.FamilyID]; ok {
		if _, ok := byUser[collab.UserID]; ok {
			cp := *collab
			f.collabs[collab.FamilyID][collab.UserID] = &cp
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeFamilyCollaboratorRepo) Delete(_ context.Context, familyID, userID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if byUser, ok := f.collabs[familyID]; ok {
		if _, ok := byUser[userID]; ok {
			delete(byUser, userID)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeFamilyCollaboratorRepo) CountByRole(_ context.Context, familyID uuid.UUID, role domain.FamilyRole) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var count int64
	for _, collab := range f.collabs[familyID] {
		if collab.Role == role {
			count++
		}
	}
	return count, nil
}

func (f *fakeFamilyCollaboratorRepo) seedOwner(familyID, userID uuid.UUID) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.collabs[familyID] == nil {
		f.collabs[familyID] = make(map[uuid.UUID]*domain.FamilyCollaborator)
	}
	f.collabs[familyID][userID] = &domain.FamilyCollaborator{
		ID:       uuid.New(),
		FamilyID: familyID,
		UserID:   userID,
		Role:     domain.FamilyRoleOwner,
	}
}

type fakePersonRepo struct {
	mu      sync.Mutex
	persons map[uuid.UUID]*domain.Person
}

func newFakePersonRepo() *fakePersonRepo {
	return &fakePersonRepo{persons: make(map[uuid.UUID]*domain.Person)}
}

func (f *fakePersonRepo) CountByFamily(_ context.Context, familyID uuid.UUID) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var count int64
	for _, p := range f.persons {
		if p.FamilyID == familyID {
			count++
		}
	}
	return count, nil
}

func (f *fakePersonRepo) Create(_ context.Context, person *domain.Person) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *person
	f.persons[person.ID] = &cp
	return nil
}

func (f *fakePersonRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Person, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if p, ok := f.persons[id]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakePersonRepo) Update(_ context.Context, person *domain.Person) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.persons[person.ID]; !ok {
		return domain.ErrNotFound
	}
	cp := *person
	f.persons[person.ID] = &cp
	return nil
}

func (f *fakePersonRepo) SoftDeleteWithRelationships(_ context.Context, personID uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.persons[personID]; !ok {
		return domain.ErrNotFound
	}
	delete(f.persons, personID)
	return nil
}

func (f *fakePersonRepo) ListByFamily(_ context.Context, familyID uuid.UUID, search string, offset, limit int) ([]domain.Person, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	items := make([]domain.Person, 0)
	for _, p := range f.persons {
		if p.FamilyID == familyID {
			items = append(items, *p)
		}
	}
	total := int64(len(items))
	if offset >= len(items) {
		return []domain.Person{}, total, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], total, nil
}

func (f *fakePersonRepo) ListByFamilyAll(_ context.Context, familyID uuid.UUID) ([]domain.Person, error) {
	items, _, err := f.ListByFamily(context.Background(), familyID, "", 0, 1000)
	return items, err
}

func (f *fakePersonRepo) ListAncestors(_ context.Context, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	return nil, nil
}

func (f *fakePersonRepo) ListDescendants(_ context.Context, personID uuid.UUID, depth int) ([]domain.PersonWithDepth, error) {
	return nil, nil
}

type fakeRelationshipRepo struct {
	mu    sync.Mutex
	rels  map[uuid.UUID]*domain.Relationship
	byKey map[string]uuid.UUID
}

func newFakeRelationshipRepo() *fakeRelationshipRepo {
	return &fakeRelationshipRepo{
		rels:  make(map[uuid.UUID]*domain.Relationship),
		byKey: make(map[string]uuid.UUID),
	}
}

func relKey(from, to uuid.UUID, relType domain.RelationshipType) string {
	return from.String() + ":" + to.String() + ":" + string(relType)
}

func (f *fakeRelationshipRepo) Create(_ context.Context, rel *domain.Relationship) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := relKey(rel.FromPersonID, rel.ToPersonID, rel.Type)
	if _, exists := f.byKey[key]; exists {
		return domain.ErrRelationshipExists
	}
	cp := *rel
	f.rels[rel.ID] = &cp
	f.byKey[key] = rel.ID
	return nil
}

func (f *fakeRelationshipRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Relationship, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if rel, ok := f.rels[id]; ok {
		cp := *rel
		return &cp, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeRelationshipRepo) Delete(_ context.Context, id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	rel, ok := f.rels[id]
	if !ok {
		return domain.ErrNotFound
	}
	delete(f.byKey, relKey(rel.FromPersonID, rel.ToPersonID, rel.Type))
	delete(f.rels, id)
	return nil
}

func (f *fakeRelationshipRepo) ListByFamily(_ context.Context, familyID uuid.UUID) ([]domain.Relationship, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]domain.Relationship, 0)
	for _, rel := range f.rels {
		if rel.FamilyID == familyID {
			out = append(out, *rel)
		}
	}
	return out, nil
}

func (f *fakeRelationshipRepo) ListForPerson(_ context.Context, personID uuid.UUID) ([]domain.Relationship, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]domain.Relationship, 0)
	for _, rel := range f.rels {
		if rel.FromPersonID == personID || rel.ToPersonID == personID {
			out = append(out, *rel)
		}
	}
	return out, nil
}

func (f *fakeRelationshipRepo) CountParentsOf(_ context.Context, childID uuid.UUID) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var count int64
	for _, rel := range f.rels {
		if rel.ToPersonID == childID && rel.Type == domain.RelationshipParent {
			count++
		}
	}
	return count, nil
}

func (f *fakeRelationshipRepo) HasAncestorPath(_ context.Context, startID, targetID uuid.UUID) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	visited := map[uuid.UUID]struct{}{}
	var walk func(uuid.UUID) bool
	walk = func(current uuid.UUID) bool {
		if current == targetID {
			return true
		}
		if _, ok := visited[current]; ok {
			return false
		}
		visited[current] = struct{}{}
		for _, rel := range f.rels {
			if rel.ToPersonID == current && rel.Type == domain.RelationshipParent {
				if walk(rel.FromPersonID) {
					return true
				}
			}
		}
		return false
	}
	return walk(startID), nil
}

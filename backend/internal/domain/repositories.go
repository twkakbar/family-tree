package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository persists User entities.
// Implementations translate driver errors into domain errors
// (ErrNotFound, ErrEmailTaken) and must respect soft deletion.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	List(ctx context.Context, search string, offset, limit int) ([]User, int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// FamilyRepository persists Family entities.
type FamilyRepository interface {
	Create(ctx context.Context, family *Family) error
	// CreateWithOwner atomically inserts the family and the owner's collaborator row.
	CreateWithOwner(ctx context.Context, family *Family, ownerID uuid.UUID, role FamilyRole) error
	GetByID(ctx context.Context, id uuid.UUID) (*Family, error)
	// GetByName returns a non-deleted family with the exact name.
	GetByName(ctx context.Context, name string) (*Family, error)
	Update(ctx context.Context, family *Family) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	// ListForUser returns families the user collaborates on, newest first.
	ListForUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]FamilySummary, int64, error)
	// GetSummaryForUser returns one family when the caller is a collaborator.
	GetSummaryForUser(ctx context.Context, userID, familyID uuid.UUID) (*FamilySummary, error)
}

// FamilyCollaboratorRepository persists per-family access grants.
type FamilyCollaboratorRepository interface {
	Create(ctx context.Context, collab *FamilyCollaborator) error
	Get(ctx context.Context, familyID, userID uuid.UUID) (*FamilyCollaborator, error)
	ListByFamily(ctx context.Context, familyID uuid.UUID) ([]FamilyCollaborator, error)
	Update(ctx context.Context, collab *FamilyCollaborator) error
	Delete(ctx context.Context, familyID, userID uuid.UUID) error
	CountByRole(ctx context.Context, familyID uuid.UUID, role FamilyRole) (int64, error)
}

// PersonRepository persists Person entities.
type PersonRepository interface {
	CountByFamily(ctx context.Context, familyID uuid.UUID) (int64, error)
	Create(ctx context.Context, person *Person) error
	GetByID(ctx context.Context, id uuid.UUID) (*Person, error)
	Update(ctx context.Context, person *Person) error
	// SoftDeleteWithRelationships soft-deletes the person and hard-deletes its edges.
	SoftDeleteWithRelationships(ctx context.Context, personID uuid.UUID) error
	ListByFamily(ctx context.Context, familyID uuid.UUID, search string, offset, limit int) ([]Person, int64, error)
	ListByFamilyAll(ctx context.Context, familyID uuid.UUID) ([]Person, error)
	ListAncestors(ctx context.Context, personID uuid.UUID, depth int) ([]PersonWithDepth, error)
	ListDescendants(ctx context.Context, personID uuid.UUID, depth int) ([]PersonWithDepth, error)
}

// RelationshipRepository persists relationship edges.
type RelationshipRepository interface {
	Create(ctx context.Context, rel *Relationship) error
	GetByID(ctx context.Context, id uuid.UUID) (*Relationship, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByFamily(ctx context.Context, familyID uuid.UUID) ([]Relationship, error)
	ListForPerson(ctx context.Context, personID uuid.UUID) ([]Relationship, error)
	CountParentsOf(ctx context.Context, childID uuid.UUID) (int64, error)
	// HasAncestorPath reports whether target is reachable from start by following
	// parent edges upward (start is the child, walking to parents).
	HasAncestorPath(ctx context.Context, startID, targetID uuid.UUID) (bool, error)
}

// EditAccessRequestRepository persists edit-access requests.
type EditAccessRequestRepository interface {
	Create(ctx context.Context, req *EditAccessRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*EditAccessRequest, error)
	GetLatestForUser(ctx context.Context, userID uuid.UUID) (*EditAccessRequest, error)
	ListByStatus(ctx context.Context, status AccessRequestStatus) ([]EditAccessRequest, error)
	Update(ctx context.Context, req *EditAccessRequest) error
}

// RefreshTokenRepository persists refresh token hashes.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	// GetByHash returns the token row for a SHA-256 hash regardless of its
	// revocation state; callers decide how to treat revoked tokens.
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

// PasswordResetTokenRepository persists password-reset token hashes.
type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) error
	GetByHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}

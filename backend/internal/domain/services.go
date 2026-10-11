package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RegisterInput carries validated data for creating an account.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// LoginInput carries credentials for authentication.
type LoginInput struct {
	Email    string
	Password string
}

// RegisterResult is the outcome of creating an account.
// Session is set only when email verification is not required (SMTP off).
type RegisterResult struct {
	Session              *AuthResult
	VerificationRequired bool
}

// Mailer delivers transactional email. A nil or disabled mailer means
// registration skips verification and signs the user in immediately.
type Mailer interface {
	Enabled() bool
	Send(ctx context.Context, to, subject, text string) error
}

// AuthResult is returned by every operation that establishes a session.
// The raw RefreshToken is delivered to the client exactly once (cookie);
// only its hash is stored.
type AuthResult struct {
	User                  *User
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

// AuthService implements registration and the token lifecycle
// (login, refresh rotation with reuse detection, logout).
type AuthService interface {
	Register(ctx context.Context, in RegisterInput) (*RegisterResult, error)
	Login(ctx context.Context, in LoginInput) (*AuthResult, error)
	// Refresh rotates a valid refresh token into a new token pair. Presenting
	// a revoked token revokes every session of that user (reuse detection).
	Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	// Logout revokes the presented refresh token. It is idempotent.
	Logout(ctx context.Context, refreshToken string) error
	// EnsureAdmin creates the bootstrap admin account when credentials are
	// configured and no account with that email exists. No-op otherwise.
	EnsureAdmin(ctx context.Context, name, email, password string) error
	// RequestPasswordReset creates a reset token for the email when the account
	// exists. Always succeeds from the caller's perspective (no enumeration).
	// The raw token is returned only so non-prod / admin flows can surface it;
	// production forgot-password handlers must not expose it to the client.
	RequestPasswordReset(ctx context.Context, email string) (rawToken string, err error)
	// ResetPassword consumes a valid raw reset token and sets a new password.
	ResetPassword(ctx context.Context, rawToken, newPassword string) error
	// ResendVerification emails a new link when the account exists and is
	// still unverified. Unknown or already-verified emails are a no-op.
	ResendVerification(ctx context.Context, email string) error
	// VerifyEmail consumes a one-time token and marks the account verified.
	VerifyEmail(ctx context.Context, rawToken string) error
}

// UpdateProfileInput carries the mutable fields of the caller's own profile.
type UpdateProfileInput struct {
	Name string
}

// UserService exposes account profile operations.
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, in UpdateProfileInput) (*User, error)
}

// CreateFamilyInput carries validated data for creating a family.
type CreateFamilyInput struct {
	Name        string
	Description string
}

// UpdateFamilyInput carries mutable family fields.
type UpdateFamilyInput struct {
	Name        string
	Description string
}

// AddCollaboratorInput invites an existing user to a family.
type AddCollaboratorInput struct {
	Email string
	Role  FamilyRole
}

// FamilyService implements family CRUD and collaborator management.
type FamilyService interface {
	Create(ctx context.Context, userID uuid.UUID, in CreateFamilyInput) (*FamilySummary, error)
	ListForUser(ctx context.Context, userID uuid.UUID, page, perPage int) ([]FamilySummary, int64, error)
	GetByID(ctx context.Context, userID, familyID uuid.UUID) (*FamilySummary, error)
	Update(ctx context.Context, userID, familyID uuid.UUID, in UpdateFamilyInput) (*FamilySummary, error)
	Delete(ctx context.Context, userID, familyID uuid.UUID) error
	ListCollaborators(ctx context.Context, userID, familyID uuid.UUID) ([]CollaboratorSummary, error)
	AddCollaborator(ctx context.Context, userID, familyID uuid.UUID, in AddCollaboratorInput) (*CollaboratorSummary, error)
	UpdateCollaboratorRole(ctx context.Context, userID, familyID, targetUserID uuid.UUID, role FamilyRole) (*CollaboratorSummary, error)
	RemoveCollaborator(ctx context.Context, userID, familyID, targetUserID uuid.UUID) error
}

// CreatePersonInput carries validated data for creating a person.
type CreatePersonInput struct {
	FirstName string
	LastName  string
	Gender    Gender
	BirthDate *time.Time
	DeathDate *time.Time
	Bio       string
}

// UpdatePersonInput carries mutable person fields.
type UpdatePersonInput struct {
	FirstName string
	LastName  string
	Gender    Gender
	BirthDate *time.Time
	DeathDate *time.Time
	Bio       string
}

// CreateRelationshipInput carries validated data for a new edge.
type CreateRelationshipInput struct {
	FromPersonID uuid.UUID
	ToPersonID   uuid.UUID
	Type         RelationshipType
}

// PersonService implements person CRUD, tree data, and relationship management.
type PersonService interface {
	ListByFamily(ctx context.Context, userID, familyID uuid.UUID, search string, page, perPage int) ([]Person, int64, error)
	Create(ctx context.Context, userID, familyID uuid.UUID, in CreatePersonInput) (*Person, error)
	GetByID(ctx context.Context, userID, personID uuid.UUID) (*PersonDetail, error)
	Update(ctx context.Context, userID, personID uuid.UUID, in UpdatePersonInput) (*Person, error)
	Delete(ctx context.Context, userID, personID uuid.UUID) error
	GetFamilyTree(ctx context.Context, userID, familyID uuid.UUID) (*FamilyTree, error)
	GetAncestors(ctx context.Context, userID, personID uuid.UUID, depth int) ([]PersonWithDepth, error)
	GetDescendants(ctx context.Context, userID, personID uuid.UUID, depth int) ([]PersonWithDepth, error)
	CreateRelationship(ctx context.Context, userID, familyID uuid.UUID, in CreateRelationshipInput) (*Relationship, error)
	DeleteRelationship(ctx context.Context, userID, relationshipID uuid.UUID) error
	UpdatePhoto(ctx context.Context, userID, personID uuid.UUID, data []byte) (*Person, error)
}

// PublicTreeService exposes unauthenticated reads of the royal family tree.
type PublicTreeService interface {
	Meta(ctx context.Context) (*Family, error)
	Tree(ctx context.Context) (*FamilyTree, error)
	GetPerson(ctx context.Context, personID uuid.UUID) (*PersonDetail, error)
	Ancestors(ctx context.Context, personID uuid.UUID, depth int) ([]PersonWithDepth, error)
	Descendants(ctx context.Context, personID uuid.UUID, depth int) ([]PersonWithDepth, error)
}

// AccessRequestService manages edit-access requests for the royal family.
type AccessRequestService interface {
	Submit(ctx context.Context, userID uuid.UUID, message string) (*EditAccessRequest, error)
	MyLatest(ctx context.Context, userID uuid.UUID) (*EditAccessRequest, error)
	ListPending(ctx context.Context, adminID uuid.UUID) ([]AccessRequestSummary, error)
	Approve(ctx context.Context, adminID, requestID uuid.UUID) (*EditAccessRequest, error)
	Reject(ctx context.Context, adminID, requestID uuid.UUID) (*EditAccessRequest, error)
}

// RoyalFamilyService seeds and resolves the singleton public family.
type RoyalFamilyService interface {
	// Ensure creates the royal family owned by ownerID when missing.
	Ensure(ctx context.Context, ownerID uuid.UUID, name, description string) (*Family, error)
	// ID returns the seeded royal family id (ErrNotFound if not seeded).
	ID(ctx context.Context) (uuid.UUID, error)
	Get(ctx context.Context) (*Family, error)
	// CanEdit reports whether the user may mutate the royal tree.
	CanEdit(ctx context.Context, userID uuid.UUID) (bool, error)
}

// TransferService exports and imports the family tree as portable JSON.
type TransferService interface {
	ExportJSON(ctx context.Context, userID, familyID uuid.UUID) (*TreeTransfer, error)
	ImportJSON(ctx context.Context, userID, familyID uuid.UUID, doc *TreeTransfer) (*ImportResult, error)
}

// AdminUpdateUserInput is the admin patch for an account.
type AdminUpdateUserInput struct {
	Name          *string
	Role          *Role
	EmailVerified *bool
}

// AdminUserService is super-admin account management.
type AdminUserService interface {
	List(ctx context.Context, adminID uuid.UUID, search string, page, perPage int) ([]User, int64, error)
	Update(ctx context.Context, adminID, userID uuid.UUID, in AdminUpdateUserInput) (*User, error)
	Delete(ctx context.Context, adminID, userID uuid.UUID) error
	// IssuePasswordReset creates a reset token and returns the raw token once
	// for the admin to share out-of-band (no SMTP yet).
	IssuePasswordReset(ctx context.Context, adminID, userID uuid.UUID) (rawToken string, err error)
}

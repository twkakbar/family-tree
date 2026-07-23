package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Family is a container for one family tree, owned by the creating user and
// optionally shared with collaborators.
type Family struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null;default:''"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// FamilyRole is a user's permission level within a single family.
type FamilyRole string

// Per-family roles, from most to least privileged.
const (
	FamilyRoleOwner  FamilyRole = "owner"
	FamilyRoleEditor FamilyRole = "editor"
	FamilyRoleViewer FamilyRole = "viewer"
)

// FamilyCollaborator grants a user access to a family with a specific role.
// Grants are hard-deleted when access is revoked.
type FamilyCollaborator struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	FamilyID  uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_family_user"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_family_user"`
	Role      FamilyRole `gorm:"type:text;not null;default:viewer"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CanEdit reports whether the collaborator may modify family content.
func (c *FamilyCollaborator) CanEdit() bool {
	return c.Role == FamilyRoleOwner || c.Role == FamilyRoleEditor
}

var familyRoleRank = map[FamilyRole]int{
	FamilyRoleOwner:  3,
	FamilyRoleEditor: 2,
	FamilyRoleViewer: 1,
}

// Valid reports whether the role is one of the supported family roles.
func (r FamilyRole) Valid() bool {
	_, ok := familyRoleRank[r]
	return ok
}

// AtLeast reports whether r meets or exceeds the minimum required role.
func (r FamilyRole) AtLeast(min FamilyRole) bool {
	return familyRoleRank[r] >= familyRoleRank[min]
}

// FamilySummary is a family visible to a collaborator, including their role
// and a denormalized persons count for list/detail views.
type FamilySummary struct {
	Family       Family
	Role         FamilyRole
	PersonsCount int64
}

// CollaboratorSummary joins a collaborator grant with the invited user's
// public profile fields.
type CollaboratorSummary struct {
	Collaborator FamilyCollaborator
	User         User
}

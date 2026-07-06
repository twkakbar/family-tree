package domain

import (
	"time"

	"github.com/google/uuid"
)

// RelationshipType classifies an edge between two persons.
type RelationshipType string

// Supported relationship types.
//
//   - parent: directed edge, FromPersonID is the parent, ToPersonID the child.
//   - spouse: undirected edge, stored once in canonical order
//     (FromPersonID < ToPersonID, enforced by a DB CHECK constraint).
const (
	RelationshipParent RelationshipType = "parent"
	RelationshipSpouse RelationshipType = "spouse"
)

// Relationship is a typed edge between two persons of the same family.
// Edges are hard-deleted: soft-deleted edges would silently corrupt
// ancestor/descendant traversals.
type Relationship struct {
	ID           uuid.UUID        `gorm:"type:uuid;primaryKey"`
	FamilyID     uuid.UUID        `gorm:"type:uuid;not null;index"`
	FromPersonID uuid.UUID        `gorm:"type:uuid;not null;index"`
	ToPersonID   uuid.UUID        `gorm:"type:uuid;not null;index"`
	Type         RelationshipType `gorm:"type:text;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

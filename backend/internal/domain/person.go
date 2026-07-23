package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Gender of a person record.
type Gender string

// Supported gender values.
const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderUnknown Gender = "unknown"
)

// Person is a node in a family tree (a family member). Persons are pure data
// records scoped to exactly one family; they have no login and are unrelated
// to User accounts.
type Person struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	FamilyID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	FirstName string     `gorm:"not null"`
	LastName  string     `gorm:"not null;default:''"`
	Gender    Gender     `gorm:"type:text;not null;default:unknown"`
	BirthDate *time.Time `gorm:"type:date"`
	DeathDate *time.Time `gorm:"type:date"`
	Bio       string     `gorm:"not null;default:''"`
	PhotoPath string     `gorm:"not null;default:''"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName overrides GORM's default pluralization ("people").
func (Person) TableName() string { return "persons" }

// FullName returns the display name of the person.
func (p *Person) FullName() string {
	if p.LastName == "" {
		return p.FirstName
	}
	return p.FirstName + " " + p.LastName
}

// Valid reports whether the gender value is supported.
func (g Gender) Valid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther, GenderUnknown:
		return true
	default:
		return false
	}
}

// PersonWithDepth is a person returned by ancestor/descendant traversals.
type PersonWithDepth struct {
	Person Person
	Depth  int
}

// PersonDetail is a person with their incident relationship edges.
type PersonDetail struct {
	Person        Person
	Relationships []Relationship
}

// FamilyTree is the full graph payload for client-side rendering.
type FamilyTree struct {
	Persons       []Person
	Relationships []Relationship
}

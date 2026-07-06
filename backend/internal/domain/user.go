package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role is the global application role of a user account.
type Role string

// Global roles. Per-family permissions are modeled by FamilyRole instead.
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User is an account that can authenticate against the API.
// A User is NOT a family member; tree nodes are modeled by Person.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"type:citext;not null"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         Role      `gorm:"type:text;not null;default:user"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// IsAdmin reports whether the user has the global admin role.
func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }

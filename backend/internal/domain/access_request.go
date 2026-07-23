package domain

import (
	"time"

	"github.com/google/uuid"
)

// AccessRequestStatus is the lifecycle of an edit-access request.
type AccessRequestStatus string

const (
	AccessRequestPending  AccessRequestStatus = "pending"
	AccessRequestApproved AccessRequestStatus = "approved"
	AccessRequestRejected AccessRequestStatus = "rejected"
)

// EditAccessRequest is a registered user asking to edit the royal family tree.
type EditAccessRequest struct {
	ID         uuid.UUID           `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID           `gorm:"type:uuid;not null;index"`
	Message    string              `gorm:"not null;default:''"`
	Status     AccessRequestStatus `gorm:"type:text;not null;default:pending"`
	ReviewedBy *uuid.UUID          `gorm:"type:uuid"`
	ReviewedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (EditAccessRequest) TableName() string { return "edit_access_requests" }

// AccessRequestSummary joins a request with the applicant's public profile.
type AccessRequestSummary struct {
	Request EditAccessRequest
	User    User
}

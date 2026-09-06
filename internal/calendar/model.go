package calendar

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusActive    Status = "ACTIVE"
	StatusCancelled Status = "CANCELLED"
)

type ExternalEvent struct {
	ExternalID      string
	ExternalVersion string

	Title       string
	Description string
	Location    string

	StartsAt time.Time
	EndsAt   time.Time

	AllDay     bool
	RawPayload json.RawMessage
}

type Event struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	WorkspaceID uuid.UUID `gorm:"type:uuid;not null;index"`
	SourceID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_source_external"`

	ExternalID      string `gorm:"not null;uniqueIndex:idx_source_external"`
	ExternalVersion string

	Title       string `gorm:"not null"`
	Description string
	Location    string

	StartsAt time.Time `gorm:"not null;index"`
	EndsAt   time.Time `gorm:"not null"`

	AllDay bool

	Status Status `gorm:"type:varchar(32);not null;default:'ACTIVE'"`

	RawPayload json.RawMessage `gorm:"type:jsonb"`

	FirstSeenAt time.Time
	LastSeenAt  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

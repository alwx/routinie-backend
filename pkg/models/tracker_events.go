package models

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

var TrackerEventTypes = []string{
	"set",
	"timer_start",
	"timer_stop",
}
var DefaultTrackerEventType = "set"

var ErrTooLateToAddTrackerEvent = errors.New("too late to add tracker event")

type NewTrackerEvent struct {
	UserID    uuid.UUID `json:"-"`
	TrackerID uuid.UUID `json:"tracker_id" binding:"required"`

	Type           string         `json:"type" binding:"required"`
	Value          int            `json:"value" binding:"required"`
	AssignedToDate datatypes.Date `json:"assigned_to_date" binding:"required"`
}

type PatchedTrackerEvent struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`

	Type           *string         `json:"type"`
	Value          *int            `json:"value"`
	AssignedToDate *datatypes.Date `json:"assigned_to_date" binding:"required"`
}

type DeletedTrackerEvent struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
}

type TrackerEvent struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key"`
	UserID    *uuid.UUID `json:"user_id,omitempty" gorm:"foreignkey:id"`
	TrackerID uuid.UUID  `json:"tracker_id" gorm:"foreignkey:id"`

	Type             *string         `json:"type" gorm:"not null;default:set"`
	Value            *int            `json:"value"`
	AssignedToDate   *datatypes.Date `json:"assigned_to_date"`
	TrackerFulfilled *bool           `json:"tracker_fulfilled"`
	Metadata         datatypes.JSON  `json:"metadata,omitempty" gorm:"type:JSONB" sql:"DEFAULT:'{}'"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

type TrackerEventRepository interface {
	FindByID(transaction Transaction, id uuid.UUID) (*TrackerEvent, error)
	FindAllForUserID(userID uuid.UUID, since int, until int, isPublic bool) ([]TrackerEvent, error)

	Insert(transaction Transaction, user User, existingTracker Tracker, newTrackerEvent NewTrackerEvent) (*TrackerEvent, error)
	Patch(transaction Transaction, user User, existingTracker Tracker, patchedTrackerEvent PatchedTrackerEvent) error
	Delete(transaction Transaction, deletedTrackerEvent DeletedTrackerEvent) error
}

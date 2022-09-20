package models

import (
	"time"

	"github.com/google/uuid"
)

var TrackerTypes = []string{
	"daily",
	"timer",
}
var DefaultTrackerType = "daily"

type NewTracker struct {
	UserID          uuid.UUID  `json:"-"`
	ParentTrackerID *uuid.UUID `json:"parent_tracker_id"`

	Title string `json:"title" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Color string `json:"color" binding:"required"`

	DefaultValue  int    `json:"default_value" binding:"required"`
	GoalValue     int    `json:"goal_value" binding:"required"`
	Measurement   string `json:"measurement"`
	DefaultChange int    `json:"default_change" binding:"required"`
	IsInfinite    bool   `json:"is_infinite"`
	IsPublic      bool   `json:"is_public"`
	Rank          string `json:"rank" binding:"required"`
}

type PatchedTracker struct {
	ID              uuid.UUID  `json:"-"`
	UserID          uuid.UUID  `json:"-"`
	ParentTrackerID *uuid.UUID `json:"parent_tracker_id"`

	Title *string `json:"title"`
	Type  *string `json:"type"`
	Color *string `json:"color"`

	DefaultValue  *int    `json:"default_value"`
	GoalValue     *int    `json:"goal_value"`
	Measurement   *string `json:"measurement"`
	DefaultChange *int    `json:"default_change"`
	IsInfinite    *bool   `json:"is_infinite"`
	IsPublic      *bool   `json:"is_public"`
	Rank          *string `json:"rank"`
}

type DeletedTracker struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
}

type Tracker struct {
	ID              uuid.UUID  `json:"id" gorm:"primary_key"`
	UserID          *uuid.UUID `json:"user_id,omitempty" gorm:"foreignkey:id"`
	ParentTrackerID *uuid.UUID `json:"parent_tracker_id,omitempty" gorm:"foreignkey:id"`

	Title *string `json:"title" gorm:"not null"`
	Type  *string `json:"type" gorm:"not null"`
	Color *string `json:"color" gorm:"not null"`

	DefaultValue  *int    `json:"default_value" gorm:"not null"`
	GoalValue     *int    `json:"goal_value" gorm:"not null"`
	Measurement   *string `json:"measurement,omitempty"`
	DefaultChange *int    `json:"default_change" gorm:"not null"`
	IsInfinite    *bool   `json:"is_infinite" gorm:"default:false"`
	IsPublic      *bool   `json:"is_public,omitempty" gorm:"default:false"`
	Rank          *string `json:"rank" sql:"index" gorm:"not null"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`

	MaxStreak     int `json:"max_streak,omitempty" gorm:"->"`
	CurrentStreak int `json:"current_streak,omitempty" gorm:"->"`
}

type TrackerRepository interface {
	FindByID(transaction Transaction, id uuid.UUID) (*Tracker, error)
	FindByIDWithStreaks(transaction Transaction, id uuid.UUID, timezoneOffset int) (*Tracker, error)

	FindAllForUserID(userID uuid.UUID, since int) ([]Tracker, error)
	FindAllForUserIDWithStreaks(userID uuid.UUID, since int, timezoneOffset int) ([]Tracker, error)
	FindAllPublicForUserIDWithStreaks(userID uuid.UUID, since int, timezoneOffset int) ([]Tracker, error)

	Insert(transaction Transaction, newTracker NewTracker) (*Tracker, error)
	Patch(transaction Transaction, patchedTracker PatchedTracker) error
	Delete(transaction Transaction, deletedTracker DeletedTracker) error
}

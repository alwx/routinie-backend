package models

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type SampleTracker struct {
	ID          string         `json:"id" gorm:"primary_key"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Emoji       string         `json:"emoji" gorm:"not null"`
	Tags        pq.StringArray `json:"tags" gorm:"type:varchar(64)[]" sql:"index"`
	Data        datatypes.JSON `json:"data" gorm:"type:JSONB" sql:"DEFAULT:'{}'"`
	Priority    int            `json:"priority"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

type SampleTrackerRepository interface {
	FindAllBySimilarTag(tag string, limit int) ([]SampleTracker, error)
}

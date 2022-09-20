package utils

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func NewNullTime() *NullTime {
	var now NullTime
	now.Time = time.Now()
	now.Valid = true
	return &now
}

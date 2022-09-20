package templates

import (
	"fmt"
	"time"
)

func SafeBool(s *bool) bool {
	if s == nil {
		return false
	}
	return *s
}

func withDefault(s string, def string) string {
	if s != "" {
		return s
	}
	return def
}

func pageTitle(s string) string {
	ending := "Routinie"
	if s == ending {
		return s
	}
	return fmt.Sprintf("%s · %s", s, ending)
}

func formatDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

func humanReadableDate(t time.Time) string {
	return fmt.Sprintf("%d/%02d/%d", t.Day(), t.Month(), t.Year())
}

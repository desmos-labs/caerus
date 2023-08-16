package testutils

import (
	"time"
)

// GetTimePointer returns a pointer to the given time
func GetTimePointer(t time.Time) *time.Time {
	return &t
}

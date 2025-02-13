package domain

import (
	"time"

	"github.com/google/uuid"
)

// GetTimeNow is a function that returns the current time.
// It can be overridden in tests to provide deterministic behavior.
var GetTimeNow = func() time.Time {
	return time.Now()
}

// GetUUID is a function that generates a UUID.
// It can be overridden in tests to provide deterministic behavior.
var GetUUID = func() string {
	return uuid.NewString()
}

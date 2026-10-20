package hospital

import "errors"

// Sentinel errors of the package. Callers compare against them with
// errors.Is, because methods wrap them with extra context using %w.
var (
	// ErrInvalidData is returned when a constructor or method receives
	// values that break its rules (empty text, non-positive numbers, nil).
	ErrInvalidData = errors.New("invalid data")
)

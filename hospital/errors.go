package hospital

import "errors"

// Sentinel errors of the package. Callers compare against them with
// errors.Is, because methods wrap them with extra context using %w.
var (
	// ErrInvalidData is returned when a constructor or method receives
	// values that break its rules (empty text, non-positive numbers, nil).
	ErrInvalidData = errors.New("invalid data")

	// ErrAlreadyAsleep is returned when a sleep attack hits a patient who
	// is already asleep.
	ErrAlreadyAsleep = errors.New("patient already asleep")

	// ErrAlreadyAwake is returned when a patient who is already awake is
	// asked to wake up.
	ErrAlreadyAwake = errors.New("patient already awake")

	// ErrRoomFull is returned when a patient tries to occupy a room with
	// no free beds.
	ErrRoomFull = errors.New("room is full")

	// ErrAlreadyInRoom is returned when a patient tries to occupy a room
	// they are already in.
	ErrAlreadyInRoom = errors.New("patient already in room")

	// ErrNotInRoom is returned when releasing a patient who is not in the
	// room.
	ErrNotInRoom = errors.New("patient not in room")

	// ErrAlreadyUnderCare is returned when a doctor diagnoses a patient who
	// is already in their care.
	ErrAlreadyUnderCare = errors.New("patient already under care")

	// ErrPatientAwake is returned when someone tries to attend a sleep
	// emergency of a patient who is awake.
	ErrPatientAwake = errors.New("patient is awake")

	// ErrDuplicateID is returned when the hospital registers a patient, a
	// staff member or a room whose identifier is already registered.
	ErrDuplicateID = errors.New("duplicate id")

	// ErrNotFound is returned when the hospital looks for a patient, a
	// doctor or a room that is not registered.
	ErrNotFound = errors.New("not found")
)

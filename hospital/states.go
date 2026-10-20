package hospital

import "fmt"

// PatientState tells whether a patient is awake or where they fell asleep.
type PatientState int

// Possible states of a patient. The zero value is Awake, so a patient
// starts awake unless something changes it.
const (
	Awake PatientState = iota
	AsleepInHallway
	AsleepInBed
)

// String returns a readable name for the state, so fmt prints it as text.
func (s PatientState) String() string {
	switch s {
	case Awake:
		return "Awake"
	case AsleepInHallway:
		return "Asleep in hallway"
	case AsleepInBed:
		return "Asleep in bed"
	default:
		return fmt.Sprintf("PatientState(%d)", int(s))
	}
}

// NarcolepsyLevel is how severe the narcolepsy of a patient is.
type NarcolepsyLevel int

// Possible narcolepsy levels, from the lightest to the most severe.
const (
	Mild NarcolepsyLevel = iota
	Moderate
	Severe
)

// String returns a readable name for the level, so fmt prints it as text.
func (l NarcolepsyLevel) String() string {
	switch l {
	case Mild:
		return "Mild"
	case Moderate:
		return "Moderate"
	case Severe:
		return "Severe"
	default:
		return fmt.Sprintf("NarcolepsyLevel(%d)", int(l))
	}
}

// RoomState tells whether a room still has free beds.
type RoomState int

// Possible states of a room. The zero value is Available.
const (
	Available RoomState = iota
	Occupied
)

// String returns a readable name for the state, so fmt prints it as text.
func (s RoomState) String() string {
	switch s {
	case Available:
		return "Available"
	case Occupied:
		return "Occupied"
	default:
		return fmt.Sprintf("RoomState(%d)", int(s))
	}
}

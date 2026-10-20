package hospital

import (
	"fmt"
	"slices"
)

// Room is a hospital room with a fixed number of beds. Its fields are
// unexported, so patients enter and leave only through its methods.
type Room struct {
	number    int
	capacity  int
	state     RoomState
	occupants []*Patient
}

// NewRoom validates the data and builds an empty, available Room. It
// returns an error wrapping ErrInvalidData if number or capacity are not
// positive.
func NewRoom(number, capacity int) (*Room, error) {
	if number <= 0 {
		return nil, fmt.Errorf("new room: number %d must be positive: %w", number, ErrInvalidData)
	}
	if capacity <= 0 {
		return nil, fmt.Errorf("new room %d: capacity %d must be positive: %w", number, capacity, ErrInvalidData)
	}
	return &Room{number: number, capacity: capacity, state: Available}, nil
}

// IsAvailable reports whether the room still has at least one free bed.
func (r *Room) IsAvailable() bool {
	return len(r.occupants) < r.capacity
}

// Occupy puts patient p in a bed of the room. It returns an error wrapping
// ErrInvalidData if p is nil, ErrRoomFull if there is no free bed, or
// ErrAlreadyInRoom if p is already in this room. When the last bed is
// taken, the room becomes Occupied.
func (r *Room) Occupy(p *Patient) error {
	if p == nil {
		return fmt.Errorf("room %d: nil patient: %w", r.number, ErrInvalidData)
	}
	if !r.IsAvailable() {
		return fmt.Errorf("room %d: %w", r.number, ErrRoomFull)
	}
	if slices.Contains(r.occupants, p) {
		return fmt.Errorf("room %d: patient %s: %w", r.number, p.ID(), ErrAlreadyInRoom)
	}
	r.occupants = append(r.occupants, p)
	p.moveToBed(r)
	if !r.IsAvailable() {
		r.state = Occupied
	}
	return nil
}

// Release takes patient p out of the room, which becomes Available again.
// It returns an error wrapping ErrNotInRoom if p is not in this room.
func (r *Room) Release(p *Patient) error {
	i := slices.Index(r.occupants, p)
	if i == -1 {
		return fmt.Errorf("room %d: %w", r.number, ErrNotInRoom)
	}
	r.occupants = slices.Delete(r.occupants, i, i+1)
	r.state = Available
	return nil
}

// Number returns the number that identifies the room.
func (r *Room) Number() int {
	return r.number
}

// Capacity returns how many beds the room has.
func (r *Room) Capacity() int {
	return r.capacity
}

// State returns whether the room is Available or Occupied.
func (r *Room) State() RoomState {
	return r.state
}

// Occupants returns a copy of the patients in the room, so callers cannot
// change the room by modifying the returned slice.
func (r *Room) Occupants() []*Patient {
	return slices.Clone(r.occupants)
}

package hospital

import (
	"fmt"
	"time"
)

// Orderly is a hospital orderly. It embeds Person by value, so ID, Name and
// Age are promoted to Orderly, and it counts how many episodes it attended.
type Orderly struct {
	Person
	attended int
}

// Compile-time check: the build fails if *Orderly stops implementing Attender.
var _ Attender = (*Orderly)(nil)

// NewOrderly validates the data and builds an Orderly with no episodes
// attended. It returns an error wrapping ErrInvalidData if the person data
// are invalid.
func NewOrderly(id, name string, age int) (*Orderly, error) {
	person, err := newPerson(id, name, age)
	if err != nil {
		return nil, fmt.Errorf("new orderly: %w", err)
	}
	return &Orderly{Person: person}, nil
}

// Attend handles patient p, asleep at location, and returns the record of
// the episode. The record ID uses the own counter of the orderly, as in
// "O-01-E1". An orderly moves patients but is not a physician, so it never
// diagnoses p nor assigns them a doctor. It returns an error wrapping
// ErrInvalidData if p is nil, or ErrPatientAwake if p is awake.
func (o *Orderly) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, fmt.Errorf("orderly %s: nil patient: %w", o.ID(), ErrInvalidData)
	}
	if p.State() == Awake {
		return EpisodeRecord{}, fmt.Errorf("orderly %s: patient %s: %w", o.ID(), p.ID(), ErrPatientAwake)
	}
	id := fmt.Sprintf("%s-E%d", o.ID(), o.attended+1)
	rec := newEpisodeRecord(id, time.Now(), p, o, location, p.Room())
	o.attended++
	return rec, nil
}

// Attended returns how many episodes the orderly has attended.
func (o *Orderly) Attended() int {
	return o.attended
}

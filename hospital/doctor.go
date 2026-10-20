package hospital

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// Doctor is a physician of the hospital. It embeds Person by value, so ID,
// Name and Age are promoted to Doctor, and it keeps its own patients and the
// episodes it attended, so MyEpisodes works without the hospital.
type Doctor struct {
	Person
	specialty string
	patients  []*Patient
	episodes  []EpisodeRecord
}

// Compile-time check: the build fails if *Doctor stops implementing Attender.
var _ Attender = (*Doctor)(nil)

// NewDoctor validates the data and builds a Doctor with no patients and no
// episodes. It returns an error wrapping ErrInvalidData if the person data
// are invalid or specialty is blank.
func NewDoctor(id, name string, age int, specialty string) (*Doctor, error) {
	person, err := newPerson(id, name, age)
	if err != nil {
		return nil, fmt.Errorf("new doctor: %w", err)
	}
	specialty = strings.TrimSpace(specialty)
	if specialty == "" {
		return nil, fmt.Errorf("new doctor %s: empty specialty: %w", person.ID(), ErrInvalidData)
	}
	return &Doctor{Person: person, specialty: specialty}, nil
}

// DiagnosePatient puts patient p in the care of the doctor. It returns an
// error wrapping ErrInvalidData if p is nil, or ErrAlreadyUnderCare if p is
// already one of the patients of the doctor.
func (d *Doctor) DiagnosePatient(p *Patient) error {
	if p == nil {
		return fmt.Errorf("doctor %s: nil patient: %w", d.ID(), ErrInvalidData)
	}
	if slices.Contains(d.patients, p) {
		return fmt.Errorf("doctor %s: patient %s: %w", d.ID(), p.ID(), ErrAlreadyUnderCare)
	}
	d.patients = append(d.patients, p)
	p.assignDoctor(d)
	return nil
}

// AttendSleepEmergency checks that patient p is asleep and, if no doctor is
// in charge of p yet, diagnoses them. It returns an error wrapping
// ErrInvalidData if p is nil, or ErrPatientAwake if p is awake.
func (d *Doctor) AttendSleepEmergency(p *Patient) error {
	if p == nil {
		return fmt.Errorf("doctor %s: nil patient: %w", d.ID(), ErrInvalidData)
	}
	if p.State() == Awake {
		return fmt.Errorf("doctor %s: patient %s: %w", d.ID(), p.ID(), ErrPatientAwake)
	}
	if p.AssignedDoctor() == nil {
		return d.DiagnosePatient(p)
	}
	return nil
}

// Attend handles patient p, asleep at location, and returns the record of
// the episode, which is also saved in the episodes of the doctor. The record
// ID uses the own counter of the doctor, as in "D-01-E1". It returns the
// error of AttendSleepEmergency if p cannot be attended.
func (d *Doctor) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if err := d.AttendSleepEmergency(p); err != nil {
		return EpisodeRecord{}, err
	}
	id := fmt.Sprintf("%s-E%d", d.ID(), len(d.episodes)+1)
	rec := newEpisodeRecord(id, time.Now(), p, d, location, p.Room())
	d.episodes = append(d.episodes, rec)
	return rec, nil
}

// MyEpisodes returns a copy of the episodes the doctor attended, so callers
// cannot change the history of the doctor by modifying the returned slice.
func (d *Doctor) MyEpisodes() []EpisodeRecord {
	return slices.Clone(d.episodes)
}

// Patients returns a copy of the patients in the care of the doctor.
func (d *Doctor) Patients() []*Patient {
	return slices.Clone(d.patients)
}

// Specialty returns the medical specialty of the doctor.
func (d *Doctor) Specialty() string {
	return d.specialty
}

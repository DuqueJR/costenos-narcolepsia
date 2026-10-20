package hospital

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// Hospital orchestrates patients, rooms and staff. Its fields are
// unexported, so everything is registered and read through its methods.
// The staff lives in a []Attender, so the hospital can dispatch any kind of
// staff member (doctor, orderly...) without knowing its concrete type.
type Hospital struct {
	name     string
	patients []*Patient
	doctors  []*Doctor
	staff    []Attender
	rooms    []*Room
	history  []EpisodeRecord
	next     int // index of the next Attender in the round-robin dispatch
}

// NewHospital builds an empty Hospital called name.
func NewHospital(name string) *Hospital {
	return &Hospital{name: name}
}

// AdmitPatient registers patient p in the hospital. It returns an error
// wrapping ErrInvalidData if p is nil, or ErrDuplicateID if a patient with
// the same ID is already admitted.
func (h *Hospital) AdmitPatient(p *Patient) error {
	if p == nil {
		return fmt.Errorf("admit patient: nil patient: %w", ErrInvalidData)
	}
	if _, err := h.findPatientLocked(p.ID()); err == nil {
		return fmt.Errorf("admit patient %s: %w", p.ID(), ErrDuplicateID)
	}
	h.patients = append(h.patients, p)
	return nil
}

// HireStaff adds staff member a to the staff of the hospital. Any type
// that implements Attender can be hired, such as an Orderly. It returns an
// error wrapping ErrInvalidData if a is nil, or ErrDuplicateID if a staff
// member with the same ID is already hired.
func (h *Hospital) HireStaff(a Attender) error {
	return h.hireStaffLocked(a)
}

// hireStaffLocked is HireStaff without taking the lock.
func (h *Hospital) hireStaffLocked(a Attender) error {
	if a == nil {
		return fmt.Errorf("hire staff: nil attender: %w", ErrInvalidData)
	}
	for _, s := range h.staff {
		if s.ID() == a.ID() {
			return fmt.Errorf("hire staff %s: %w", a.ID(), ErrDuplicateID)
		}
	}
	h.staff = append(h.staff, a)
	return nil
}

// HireDoctor adds doctor d to the staff, so it can attend emergencies, and
// to the doctors of the hospital, so it can be found with FindDoctor. It
// returns an error wrapping ErrInvalidData if d is nil, or the error of
// HireStaff.
func (h *Hospital) HireDoctor(d *Doctor) error {
	if d == nil {
		return fmt.Errorf("hire doctor: nil doctor: %w", ErrInvalidData)
	}
	if err := h.hireStaffLocked(d); err != nil {
		return fmt.Errorf("hire doctor: %w", err)
	}
	h.doctors = append(h.doctors, d)
	return nil
}

// AddRoom adds room r to the hospital. It returns an error wrapping
// ErrInvalidData if r is nil, or ErrDuplicateID if a room with the same
// number already exists.
func (h *Hospital) AddRoom(r *Room) error {
	if r == nil {
		return fmt.Errorf("add room: nil room: %w", ErrInvalidData)
	}
	if _, err := h.findRoomLocked(r.Number()); err == nil {
		return fmt.Errorf("add room %d: %w", r.Number(), ErrDuplicateID)
	}
	h.rooms = append(h.rooms, r)
	return nil
}

// AssignRoom puts patient p, asleep in the hallway, in the first room that
// has a free bed, following the order in which the rooms were added. It
// returns an error wrapping ErrInvalidData if p is nil, ErrNotInHallway if p
// is not asleep in the hallway, or ErrNoRoomAvailable if every room is full.
// In that last case the patient stays asleep in the hallway. AssignRoom
// never creates rooms.
func (h *Hospital) AssignRoom(p *Patient) (*Room, error) {
	return h.assignRoomLocked(p)
}

// assignRoomLocked is AssignRoom without taking the lock.
func (h *Hospital) assignRoomLocked(p *Patient) (*Room, error) {
	if p == nil {
		return nil, fmt.Errorf("assign room: nil patient: %w", ErrInvalidData)
	}
	if p.State() != AsleepInHallway {
		return nil, fmt.Errorf("assign room to %s: %w", p.ID(), ErrNotInHallway)
	}
	for _, r := range h.rooms {
		if r.IsAvailable() {
			if err := r.Occupy(p); err != nil {
				return nil, fmt.Errorf("assign room to %s: %w", p.ID(), err)
			}
			return r, nil
		}
	}
	return nil, fmt.Errorf("assign room to %s: %w", p.ID(), ErrNoRoomAvailable)
}

// dispatchLocked chooses the staff member who attends the next emergency,
// taking turns over the whole staff (round-robin). It returns an error
// wrapping ErrNoStaff if the hospital has no staff.
func (h *Hospital) dispatchLocked() (Attender, error) {
	if len(h.staff) == 0 {
		return nil, fmt.Errorf("dispatch: %w", ErrNoStaff)
	}
	i := h.next % len(h.staff)
	h.next = i + 1
	return h.staff[i], nil
}

// RegisterEpisode adds record rec to the history of the hospital. It
// returns an error wrapping ErrInvalidData if rec has no patient.
func (h *Hospital) RegisterEpisode(rec EpisodeRecord) error {
	return h.registerEpisodeLocked(rec)
}

// registerEpisodeLocked is RegisterEpisode without taking the lock.
func (h *Hospital) registerEpisodeLocked(rec EpisodeRecord) error {
	if rec.Patient() == nil {
		return fmt.Errorf("register episode: record without patient: %w", ErrInvalidData)
	}
	h.history = append(h.history, rec)
	return nil
}

// HandleSleepAttack runs the whole flow of a sleep attack of patient p at
// location: the patient falls asleep, gets the first free room, the next
// staff member attends them and the episode is added to the history.
//
// It returns the record AND an error at the same time. An error wrapping
// ErrNoRoomAvailable means "episode registered, but the patient stayed in
// the hallway": the returned record is valid and is already in History().
// Any other error means the episode was not registered and the record is
// empty: ErrInvalidData if p is nil, ErrAlreadyAsleep if p is already
// asleep, ErrNoStaff if there is nobody to attend, or the error of Attend.
func (h *Hospital) HandleSleepAttack(p *Patient, location string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, fmt.Errorf("handle sleep attack: nil patient: %w", ErrInvalidData)
	}
	if err := p.SufferSleepAttack(location); err != nil {
		return EpisodeRecord{}, fmt.Errorf("handle sleep attack: %w", err)
	}
	// Not having a room does not stop the flow: the patient is attended in
	// the hallway and roomErr is returned at the end.
	_, roomErr := h.assignRoomLocked(p)
	if roomErr != nil && !errors.Is(roomErr, ErrNoRoomAvailable) {
		return EpisodeRecord{}, fmt.Errorf("handle sleep attack: %w", roomErr)
	}
	attender, err := h.dispatchLocked()
	if err != nil {
		return EpisodeRecord{}, fmt.Errorf("handle sleep attack of %s: %w", p.ID(), err)
	}
	rec, err := attender.Attend(p, location)
	if err != nil {
		return EpisodeRecord{}, fmt.Errorf("handle sleep attack: %w", err)
	}
	if err := h.registerEpisodeLocked(rec); err != nil {
		return EpisodeRecord{}, fmt.Errorf("handle sleep attack: %w", err)
	}
	return rec, roomErr
}

// WakeUpPatient wakes patient p up. If p is sleeping in a room, the room is
// released first, so its bed is free for the next patient. It returns an
// error wrapping ErrInvalidData if p is nil, or the error of Release or
// WakeUp, such as ErrAlreadyAwake.
func (h *Hospital) WakeUpPatient(p *Patient) error {
	if p == nil {
		return fmt.Errorf("wake up patient: nil patient: %w", ErrInvalidData)
	}
	if r := p.Room(); r != nil {
		if err := r.Release(p); err != nil {
			return fmt.Errorf("wake up patient %s: %w", p.ID(), err)
		}
	}
	if err := p.WakeUp(); err != nil {
		return fmt.Errorf("wake up patient: %w", err)
	}
	return nil
}

// PatientsInHallway returns the admitted patients that are asleep in a
// hallway, in the order they were admitted, so an orderly can be sent to
// them. It returns an empty slice, never nil, when there are none.
func (h *Hospital) PatientsInHallway() []*Patient {
	hallway := []*Patient{}
	for _, p := range h.patients {
		if p.State() == AsleepInHallway {
			hallway = append(hallway, p)
		}
	}
	return hallway
}

// SevereReport returns, for every admitted patient with level Severe, how
// many episodes they had today according to the history of the hospital.
// Severe patients without episodes today appear with 0, and patients of
// other levels never appear. A map has no order, so the caller must sort the
// keys if it needs a stable order.
func (h *Hospital) SevereReport() map[*Patient]int {
	report := make(map[*Patient]int)
	for _, p := range h.patients {
		if p.Level() == Severe {
			report[p] = 0
		}
	}
	today := time.Now()
	for _, rec := range h.history {
		// Only patients already in the report (the Severe ones) are counted.
		if _, ok := report[rec.Patient()]; ok && sameDay(rec.At(), today) {
			report[rec.Patient()]++
		}
	}
	return report
}

// sameDay reports whether a and b fall on the same calendar day (same year,
// month and day), no matter the time of day.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// FindPatient returns the admitted patient with the given id. It returns an
// error wrapping ErrNotFound if there is none.
func (h *Hospital) FindPatient(id string) (*Patient, error) {
	return h.findPatientLocked(id)
}

// findPatientLocked is FindPatient without taking the lock.
func (h *Hospital) findPatientLocked(id string) (*Patient, error) {
	for _, p := range h.patients {
		if p.ID() == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("find patient %s: %w", id, ErrNotFound)
}

// FindDoctor returns the hired doctor with the given id. It returns an
// error wrapping ErrNotFound if there is none.
func (h *Hospital) FindDoctor(id string) (*Doctor, error) {
	for _, d := range h.doctors {
		if d.ID() == id {
			return d, nil
		}
	}
	return nil, fmt.Errorf("find doctor %s: %w", id, ErrNotFound)
}

// FindRoom returns the room with the given number. It returns an error
// wrapping ErrNotFound if there is none.
func (h *Hospital) FindRoom(number int) (*Room, error) {
	return h.findRoomLocked(number)
}

// findRoomLocked is FindRoom without taking the lock.
func (h *Hospital) findRoomLocked(number int) (*Room, error) {
	for _, r := range h.rooms {
		if r.Number() == number {
			return r, nil
		}
	}
	return nil, fmt.Errorf("find room %d: %w", number, ErrNotFound)
}

// Name returns the name of the hospital.
func (h *Hospital) Name() string {
	return h.name
}

// Patients returns a copy of the admitted patients, so callers cannot
// change the hospital by modifying the returned slice.
func (h *Hospital) Patients() []*Patient {
	return slices.Clone(h.patients)
}

// Doctors returns a copy of the hired doctors.
func (h *Hospital) Doctors() []*Doctor {
	return slices.Clone(h.doctors)
}

// Staff returns a copy of the whole staff, doctors included.
func (h *Hospital) Staff() []Attender {
	return slices.Clone(h.staff)
}

// Rooms returns a copy of the rooms of the hospital.
func (h *Hospital) Rooms() []*Room {
	return slices.Clone(h.rooms)
}

// History returns a copy of every episode registered in the hospital.
func (h *Hospital) History() []EpisodeRecord {
	return slices.Clone(h.history)
}

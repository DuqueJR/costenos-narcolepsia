package hospital

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestStatesString(t *testing.T) {
	tests := []struct {
		value fmt.Stringer
		want  string
	}{
		{Awake, "Awake"},
		{AsleepInHallway, "Asleep in hallway"},
		{AsleepInBed, "Asleep in bed"},
		{PatientState(7), "PatientState(7)"},
		{Mild, "Mild"},
		{Moderate, "Moderate"},
		{Severe, "Severe"},
		{NarcolepsyLevel(9), "NarcolepsyLevel(9)"},
		{Available, "Available"},
		{Occupied, "Occupied"},
		{RoomState(-1), "RoomState(-1)"},
	}
	for _, tt := range tests {
		if got := tt.value.String(); got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

// newTestPatient builds a valid patient or stops the test.
func newTestPatient(t *testing.T, id string) *Patient {
	t.Helper()
	p, err := NewPatient(id, "Patient "+id, 30, Moderate)
	if err != nil {
		t.Fatalf("NewPatient(%q) error: %v", id, err)
	}
	return p
}

// newTestRoom builds a valid room or stops the test.
func newTestRoom(t *testing.T, number, capacity int) *Room {
	t.Helper()
	r, err := NewRoom(number, capacity)
	if err != nil {
		t.Fatalf("NewRoom(%d, %d) error: %v", number, capacity, err)
	}
	return r
}

func TestRoomOccupySuccess(t *testing.T) {
	r := newTestRoom(t, 101, 1)
	p := newTestPatient(t, "P1")
	if err := p.SufferSleepAttack("cafeteria"); err != nil {
		t.Fatalf("SufferSleepAttack error: %v", err)
	}

	if err := r.Occupy(p); err != nil {
		t.Fatalf("Occupy error: %v", err)
	}
	if r.State() != Occupied {
		t.Errorf("room state = %v, want %v", r.State(), Occupied)
	}
	if p.State() != AsleepInBed {
		t.Errorf("patient state = %v, want %v", p.State(), AsleepInBed)
	}
	if p.Room() != r {
		t.Errorf("patient room = %v, want room %d", p.Room(), r.Number())
	}
}

func TestRoomOccupyFull(t *testing.T) {
	r := newTestRoom(t, 102, 1)
	if err := r.Occupy(newTestPatient(t, "P1")); err != nil {
		t.Fatalf("first Occupy error: %v", err)
	}

	err := r.Occupy(newTestPatient(t, "P2"))
	if !errors.Is(err, ErrRoomFull) {
		t.Errorf("Occupy on full room error = %v, want %v", err, ErrRoomFull)
	}
}

func TestRoomReleaseNotInRoom(t *testing.T) {
	r := newTestRoom(t, 103, 2)
	p := newTestPatient(t, "P1")

	err := r.Release(p)
	if !errors.Is(err, ErrNotInRoom) {
		t.Errorf("Release error = %v, want %v", err, ErrNotInRoom)
	}
}

func TestPatientWakeUpAlreadyAwake(t *testing.T) {
	p := newTestPatient(t, "P1")

	err := p.WakeUp()
	if !errors.Is(err, ErrAlreadyAwake) {
		t.Errorf("WakeUp error = %v, want %v", err, ErrAlreadyAwake)
	}
}

// fakeAttender is a minimal Attender for tests. It never declares that it
// implements Attender: having the three methods is enough.
type fakeAttender struct {
	id   string
	name string
}

func (f fakeAttender) ID() string   { return f.id }
func (f fakeAttender) Name() string { return f.name }

func (f fakeAttender) Attend(p *Patient, location string) (EpisodeRecord, error) {
	return newEpisodeRecord(f.id+"-E1", time.Now(), p, f, location, p.Room()), nil
}

func TestEpisodeRecordSummary(t *testing.T) {
	at := time.Date(2026, 10, 19, 14, 15, 0, 0, time.UTC)
	p := newTestPatient(t, "P-001")
	by := fakeAttender{id: "D-001", name: "Dr. Karen Ospina"}
	room := newTestRoom(t, 101, 2)

	tests := []struct {
		name string
		room *Room
		want string
	}{
		{"with room", room, "2026-10-19 14:15 P-001 cafeteria -> room 101 (Dr. Karen Ospina)"},
		{"hallway", nil, "2026-10-19 14:15 P-001 cafeteria -> hallway (Dr. Karen Ospina)"},
	}
	for _, tt := range tests {
		e := newEpisodeRecord("D-001-E1", at, p, by, "cafeteria", tt.room)
		if got := e.Summary(); got != tt.want {
			t.Errorf("%s: Summary() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// newTestDoctor builds a valid doctor or stops the test.
func newTestDoctor(t *testing.T, id string) *Doctor {
	t.Helper()
	d, err := NewDoctor(id, "Doctor "+id, 45, "Neurology")
	if err != nil {
		t.Fatalf("NewDoctor(%q) error: %v", id, err)
	}
	return d
}

func TestDoctorAttendAsleepPatient(t *testing.T) {
	d := newTestDoctor(t, "D-01")
	p := newTestPatient(t, "P1")
	if err := p.SufferSleepAttack("cafeteria"); err != nil {
		t.Fatalf("SufferSleepAttack error: %v", err)
	}

	rec, err := d.Attend(p, "cafeteria")
	if err != nil {
		t.Fatalf("Attend error: %v", err)
	}
	episodes := d.MyEpisodes()
	if len(episodes) != 1 {
		t.Fatalf("len(MyEpisodes()) = %d, want 1", len(episodes))
	}
	if episodes[0].ID() != "D-01-E1" || rec.ID() != "D-01-E1" {
		t.Errorf("episode ID = %q (returned %q), want %q", episodes[0].ID(), rec.ID(), "D-01-E1")
	}
	if p.AssignedDoctor() != d {
		t.Errorf("AssignedDoctor() = %v, want doctor %s", p.AssignedDoctor(), d.ID())
	}
}

func TestDoctorAttendAwakePatient(t *testing.T) {
	d := newTestDoctor(t, "D-01")
	p := newTestPatient(t, "P1")

	_, err := d.Attend(p, "cafeteria")
	if !errors.Is(err, ErrPatientAwake) {
		t.Errorf("Attend error = %v, want %v", err, ErrPatientAwake)
	}
	if n := len(d.MyEpisodes()); n != 0 {
		t.Errorf("len(MyEpisodes()) = %d, want 0", n)
	}
}

func TestDoctorDiagnoseTwice(t *testing.T) {
	d := newTestDoctor(t, "D-01")
	p := newTestPatient(t, "P1")
	if err := d.DiagnosePatient(p); err != nil {
		t.Fatalf("first DiagnosePatient error: %v", err)
	}

	err := d.DiagnosePatient(p)
	if !errors.Is(err, ErrAlreadyUnderCare) {
		t.Errorf("second DiagnosePatient error = %v, want %v", err, ErrAlreadyUnderCare)
	}
	if n := len(d.Patients()); n != 1 {
		t.Errorf("len(Patients()) = %d, want 1", n)
	}
}

// newTestOrderly builds a valid orderly or stops the test.
func newTestOrderly(t *testing.T, id string) *Orderly {
	t.Helper()
	o, err := NewOrderly(id, "Orderly "+id, 28)
	if err != nil {
		t.Fatalf("NewOrderly(%q) error: %v", id, err)
	}
	return o
}

func TestStaffAttendPolymorphism(t *testing.T) {
	doctor := newTestDoctor(t, "D-01")
	orderly := newTestOrderly(t, "O-01")
	staff := []Attender{doctor, orderly}

	patients := []*Patient{newTestPatient(t, "P1"), newTestPatient(t, "P2")}
	for _, p := range patients {
		if err := p.SufferSleepAttack("cafeteria"); err != nil {
			t.Fatalf("SufferSleepAttack(%s) error: %v", p.ID(), err)
		}
	}

	for i, a := range staff {
		rec, err := a.Attend(patients[i], "cafeteria")
		if err != nil {
			t.Fatalf("%s Attend error: %v", a.ID(), err)
		}
		if rec.AttendedBy() != a {
			t.Errorf("AttendedBy() = %s, want %s", rec.AttendedBy().ID(), a.ID())
		}
		if want := a.ID() + "-E1"; rec.ID() != want {
			t.Errorf("record ID = %q, want %q", rec.ID(), want)
		}
	}

	if orderly.Attended() != 1 {
		t.Errorf("Attended() = %d, want 1", orderly.Attended())
	}
	if patients[1].AssignedDoctor() != nil {
		t.Errorf("orderly assigned doctor %v, want nil", patients[1].AssignedDoctor())
	}
}

func TestOrderlyAttendAwakePatient(t *testing.T) {
	o := newTestOrderly(t, "O-01")
	p := newTestPatient(t, "P1")

	_, err := o.Attend(p, "cafeteria")
	if !errors.Is(err, ErrPatientAwake) {
		t.Errorf("Attend error = %v, want %v", err, ErrPatientAwake)
	}
	if o.Attended() != 0 {
		t.Errorf("Attended() = %d, want 0", o.Attended())
	}
}

func TestHospitalAdmitPatientDuplicate(t *testing.T) {
	h := NewHospital("Test Hospital")
	if err := h.AdmitPatient(newTestPatient(t, "P1")); err != nil {
		t.Fatalf("first AdmitPatient error: %v", err)
	}

	err := h.AdmitPatient(newTestPatient(t, "P1"))
	if !errors.Is(err, ErrDuplicateID) {
		t.Errorf("AdmitPatient duplicate error = %v, want %v", err, ErrDuplicateID)
	}
	if n := len(h.Patients()); n != 1 {
		t.Errorf("len(Patients()) = %d, want 1", n)
	}
}

func TestHospitalHireStaffDuplicate(t *testing.T) {
	h := NewHospital("Test Hospital")
	if err := h.HireDoctor(newTestDoctor(t, "S-01")); err != nil {
		t.Fatalf("HireDoctor error: %v", err)
	}

	// An orderly with the ID of a doctor is a duplicate: IDs are unique
	// across the whole staff.
	err := h.HireStaff(newTestOrderly(t, "S-01"))
	if !errors.Is(err, ErrDuplicateID) {
		t.Errorf("HireStaff duplicate error = %v, want %v", err, ErrDuplicateID)
	}
	if n := len(h.Staff()); n != 1 {
		t.Errorf("len(Staff()) = %d, want 1", n)
	}
}

func TestHospitalAddRoomDuplicate(t *testing.T) {
	h := NewHospital("Test Hospital")
	if err := h.AddRoom(newTestRoom(t, 101, 1)); err != nil {
		t.Fatalf("first AddRoom error: %v", err)
	}

	err := h.AddRoom(newTestRoom(t, 101, 2))
	if !errors.Is(err, ErrDuplicateID) {
		t.Errorf("AddRoom duplicate error = %v, want %v", err, ErrDuplicateID)
	}
	if n := len(h.Rooms()); n != 1 {
		t.Errorf("len(Rooms()) = %d, want 1", n)
	}
}

func TestHospitalNilData(t *testing.T) {
	h := NewHospital("Test Hospital")
	tests := []struct {
		name string
		err  error
	}{
		{"AdmitPatient", h.AdmitPatient(nil)},
		{"HireStaff", h.HireStaff(nil)},
		{"HireDoctor", h.HireDoctor(nil)},
		{"AddRoom", h.AddRoom(nil)},
	}
	for _, tt := range tests {
		if !errors.Is(tt.err, ErrInvalidData) {
			t.Errorf("%s(nil) error = %v, want %v", tt.name, tt.err, ErrInvalidData)
		}
	}
}

func TestHospitalFindPatientNotFound(t *testing.T) {
	h := NewHospital("Test Hospital")
	if err := h.AdmitPatient(newTestPatient(t, "P1")); err != nil {
		t.Fatalf("AdmitPatient error: %v", err)
	}

	p, err := h.FindPatient("P-999")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("FindPatient error = %v, want %v", err, ErrNotFound)
	}
	if p != nil {
		t.Errorf("FindPatient patient = %v, want nil", p)
	}
}

func TestHospitalHireDoctorAndFind(t *testing.T) {
	h := NewHospital("Test Hospital")
	d := newTestDoctor(t, "D-01")
	if err := h.HireDoctor(d); err != nil {
		t.Fatalf("HireDoctor error: %v", err)
	}
	if err := h.HireStaff(newTestOrderly(t, "O-01")); err != nil {
		t.Fatalf("HireStaff error: %v", err)
	}

	if got, err := h.FindDoctor("D-01"); err != nil || got != d {
		t.Errorf("FindDoctor = %v, %v, want doctor D-01", got, err)
	}
	if _, err := h.FindDoctor("O-01"); !errors.Is(err, ErrNotFound) {
		t.Errorf("FindDoctor(orderly) error = %v, want %v", err, ErrNotFound)
	}
	if n := len(h.Staff()); n != 2 {
		t.Errorf("len(Staff()) = %d, want 2", n)
	}
	if n := len(h.Doctors()); n != 1 {
		t.Errorf("len(Doctors()) = %d, want 1", n)
	}
}

func TestHospitalGettersReturnCopies(t *testing.T) {
	h := NewHospital("Test Hospital")
	if err := h.AdmitPatient(newTestPatient(t, "P1")); err != nil {
		t.Fatalf("AdmitPatient error: %v", err)
	}

	patients := h.Patients()
	patients[0] = nil
	if h.Patients()[0] == nil {
		t.Error("changing the slice from Patients() changed the hospital")
	}
}

// newTestHospitalWithRooms builds a hospital with the given rooms or stops
// the test.
func newTestHospitalWithRooms(t *testing.T, rooms ...*Room) *Hospital {
	t.Helper()
	h := NewHospital("Test Hospital")
	for _, r := range rooms {
		if err := h.AddRoom(r); err != nil {
			t.Fatalf("AddRoom(%d) error: %v", r.Number(), err)
		}
	}
	return h
}

// newAsleepPatient builds a patient asleep in the hallway or stops the test.
func newAsleepPatient(t *testing.T, id string) *Patient {
	t.Helper()
	p := newTestPatient(t, id)
	if err := p.SufferSleepAttack("cafeteria"); err != nil {
		t.Fatalf("SufferSleepAttack(%s) error: %v", id, err)
	}
	return p
}

func TestHospitalAssignRoomAvailable(t *testing.T) {
	full := newTestRoom(t, 101, 1)
	if err := full.Occupy(newAsleepPatient(t, "P0")); err != nil {
		t.Fatalf("Occupy error: %v", err)
	}
	free := newTestRoom(t, 102, 1)
	h := newTestHospitalWithRooms(t, full, free)
	p := newAsleepPatient(t, "P1")

	r, err := h.AssignRoom(p)
	if err != nil {
		t.Fatalf("AssignRoom error: %v", err)
	}
	if r != free {
		t.Fatalf("AssignRoom room = %v, want room %d", r, free.Number())
	}
	if r.State() != Occupied {
		t.Errorf("room state = %v, want %v", r.State(), Occupied)
	}
	if p.State() != AsleepInBed {
		t.Errorf("patient state = %v, want %v", p.State(), AsleepInBed)
	}
	if p.Room() != r {
		t.Errorf("patient room = %v, want room %d", p.Room(), r.Number())
	}
}

func TestHospitalAssignRoomNoneAvailable(t *testing.T) {
	full := newTestRoom(t, 101, 1)
	if err := full.Occupy(newAsleepPatient(t, "P0")); err != nil {
		t.Fatalf("Occupy error: %v", err)
	}
	h := newTestHospitalWithRooms(t, full)
	p := newAsleepPatient(t, "P1")

	r, err := h.AssignRoom(p)
	if !errors.Is(err, ErrNoRoomAvailable) {
		t.Errorf("AssignRoom error = %v, want %v", err, ErrNoRoomAvailable)
	}
	if r != nil {
		t.Errorf("AssignRoom room = %v, want nil", r)
	}
	if p.State() != AsleepInHallway {
		t.Errorf("patient state = %v, want %v", p.State(), AsleepInHallway)
	}
	if p.Room() != nil {
		t.Errorf("patient room = %v, want nil", p.Room())
	}
}

func TestHospitalDispatchRoundRobin(t *testing.T) {
	h := NewHospital("Test Hospital")
	doctor := newTestDoctor(t, "D-01")
	orderly := newTestOrderly(t, "O-01")
	if err := h.HireDoctor(doctor); err != nil {
		t.Fatalf("HireDoctor error: %v", err)
	}
	if err := h.HireStaff(orderly); err != nil {
		t.Fatalf("HireStaff error: %v", err)
	}

	want := []Attender{doctor, orderly, doctor, orderly}
	for i, w := range want {
		got, err := h.dispatch()
		if err != nil {
			t.Fatalf("dispatch #%d error: %v", i+1, err)
		}
		if got != w {
			t.Errorf("dispatch #%d = %s, want %s", i+1, got.ID(), w.ID())
		}
	}
}

func TestHospitalFlowErrors(t *testing.T) {
	h := NewHospital("Test Hospital")
	_, assignNilErr := h.AssignRoom(nil)
	_, assignAwakeErr := h.AssignRoom(newTestPatient(t, "P1"))
	_, dispatchErr := h.dispatch()

	tests := []struct {
		name string
		err  error
		want error
	}{
		{"AssignRoom(nil)", assignNilErr, ErrInvalidData},
		{"AssignRoom(awake)", assignAwakeErr, ErrNotInHallway},
		{"dispatch without staff", dispatchErr, ErrNoStaff},
		{"RegisterEpisode(empty)", h.RegisterEpisode(EpisodeRecord{}), ErrInvalidData},
	}
	for _, tt := range tests {
		if !errors.Is(tt.err, tt.want) {
			t.Errorf("%s error = %v, want %v", tt.name, tt.err, tt.want)
		}
	}
	if n := len(h.History()); n != 0 {
		t.Errorf("len(History()) = %d, want 0", n)
	}
}

func TestHospitalRegisterEpisode(t *testing.T) {
	h := NewHospital("Test Hospital")
	rec, err := newTestDoctor(t, "D-01").Attend(newAsleepPatient(t, "P1"), "cafeteria")
	if err != nil {
		t.Fatalf("Attend error: %v", err)
	}

	if err := h.RegisterEpisode(rec); err != nil {
		t.Fatalf("RegisterEpisode error: %v", err)
	}
	history := h.History()
	if len(history) != 1 || history[0].ID() != rec.ID() {
		t.Errorf("History() = %v, want only record %s", history, rec.ID())
	}
}

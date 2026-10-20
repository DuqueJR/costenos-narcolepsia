package hospital

import (
	"errors"
	"fmt"
	"maps"
	"slices"
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
		got, err := h.dispatchLocked()
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
	_, dispatchErr := h.dispatchLocked()

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

// newTestHospitalWithStaff builds a hospital with a doctor, an orderly and
// the given rooms or stops the test.
func newTestHospitalWithStaff(t *testing.T, rooms ...*Room) *Hospital {
	t.Helper()
	h := newTestHospitalWithRooms(t, rooms...)
	if err := h.HireDoctor(newTestDoctor(t, "D-01")); err != nil {
		t.Fatalf("HireDoctor error: %v", err)
	}
	if err := h.HireStaff(newTestOrderly(t, "O-01")); err != nil {
		t.Fatalf("HireStaff error: %v", err)
	}
	return h
}

// newTestHospital builds a hospital with a doctor (D-01), an orderly (O-01),
// the given number of single rooms (101, 102...) and three admitted awake
// patients (P1, P2, P3), or stops the test. It returns the patients in the
// order they were admitted.
func newTestHospital(t *testing.T, rooms int) (*Hospital, []*Patient) {
	t.Helper()
	var rs []*Room
	for i := 0; i < rooms; i++ {
		rs = append(rs, newTestRoom(t, 101+i, 1))
	}
	h := newTestHospitalWithStaff(t, rs...)
	patients := []*Patient{
		admitTestPatient(t, h, "P1", Moderate),
		admitTestPatient(t, h, "P2", Moderate),
		admitTestPatient(t, h, "P3", Moderate),
	}
	return h, patients
}

func TestAssignRoom_FreeRoom(t *testing.T) {
	h, patients := newTestHospital(t, 2)
	first, second := patients[0], patients[1]
	full, free := h.Rooms()[0], h.Rooms()[1]

	// first takes room 101, so second must skip it and get room 102.
	if err := first.SufferSleepAttack("cafeteria"); err != nil {
		t.Fatalf("SufferSleepAttack(%s) error: %v", first.ID(), err)
	}
	if r, err := h.AssignRoom(first); err != nil || r != full {
		t.Fatalf("AssignRoom(%s) = room %v, error %v, want room %d", first.ID(), r, err, full.Number())
	}
	if err := second.SufferSleepAttack("library"); err != nil {
		t.Fatalf("SufferSleepAttack(%s) error: %v", second.ID(), err)
	}

	r, err := h.AssignRoom(second)
	if err != nil {
		t.Fatalf("AssignRoom(%s) error: %v", second.ID(), err)
	}
	if r != free {
		t.Fatalf("AssignRoom(%s) room = %v, want room %d", second.ID(), r, free.Number())
	}
	if r.State() != Occupied {
		t.Errorf("room state = %v, want %v", r.State(), Occupied)
	}
	if second.State() != AsleepInBed {
		t.Errorf("patient state = %v, want %v", second.State(), AsleepInBed)
	}
	if second.Room() != r {
		t.Errorf("patient room = %v, want room %d", second.Room(), r.Number())
	}
	if n := len(full.Occupants()); n != 1 {
		t.Errorf("room %d has %d occupants, want 1", full.Number(), n)
	}
}

func TestAssignRoom_NoRoomAvailable(t *testing.T) {
	tests := []struct {
		name  string
		rooms int
	}{
		{"all rooms full", 1},
		{"hospital without rooms", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, patients := newTestHospital(t, tt.rooms)
			// Fill every room, so the last patient has nowhere to go.
			for _, p := range patients[:tt.rooms] {
				if err := p.SufferSleepAttack("cafeteria"); err != nil {
					t.Fatalf("SufferSleepAttack(%s) error: %v", p.ID(), err)
				}
				if _, err := h.AssignRoom(p); err != nil {
					t.Fatalf("AssignRoom(%s) error: %v", p.ID(), err)
				}
			}
			p := patients[2]
			if err := p.SufferSleepAttack("library"); err != nil {
				t.Fatalf("SufferSleepAttack(%s) error: %v", p.ID(), err)
			}

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
			if n := len(h.Rooms()); n != tt.rooms {
				t.Errorf("len(Rooms()) = %d, want %d: AssignRoom must not create rooms", n, tt.rooms)
			}
		})
	}
}

func TestHospitalHandleSleepAttackFillsRooms(t *testing.T) {
	h := newTestHospitalWithStaff(t,
		newTestRoom(t, 101, 1), newTestRoom(t, 102, 1), newTestRoom(t, 103, 1))
	patients := []*Patient{
		newTestPatient(t, "P1"), newTestPatient(t, "P2"),
		newTestPatient(t, "P3"), newTestPatient(t, "P4"),
	}

	for _, p := range patients[:3] {
		rec, err := h.HandleSleepAttack(p, "cafeteria")
		if err != nil {
			t.Fatalf("HandleSleepAttack(%s) error: %v", p.ID(), err)
		}
		if rec.Room() == nil || rec.Room() != p.Room() {
			t.Errorf("%s record room = %v, want the room of the patient", p.ID(), rec.Room())
		}
	}

	last := patients[3]
	rec, err := h.HandleSleepAttack(last, "cafeteria")
	if !errors.Is(err, ErrNoRoomAvailable) {
		t.Errorf("4th HandleSleepAttack error = %v, want %v", err, ErrNoRoomAvailable)
	}
	if rec.Patient() != last || rec.Room() != nil {
		t.Errorf("4th record = patient %v, room %v, want patient %s in the hallway", rec.Patient(), rec.Room(), last.ID())
	}
	if last.State() != AsleepInHallway {
		t.Errorf("4th patient state = %v, want %v", last.State(), AsleepInHallway)
	}
	if n := len(h.History()); n != 4 {
		t.Errorf("len(History()) = %d, want 4", n)
	}
}

func TestHospitalWakeUpPatientFreesRoom(t *testing.T) {
	room := newTestRoom(t, 101, 1)
	h := newTestHospitalWithStaff(t, room)
	inBed := newTestPatient(t, "P1")
	inHallway := newTestPatient(t, "P2")
	if _, err := h.HandleSleepAttack(inBed, "cafeteria"); err != nil {
		t.Fatalf("HandleSleepAttack(P1) error: %v", err)
	}
	if _, err := h.HandleSleepAttack(inHallway, "library"); !errors.Is(err, ErrNoRoomAvailable) {
		t.Fatalf("HandleSleepAttack(P2) error = %v, want %v", err, ErrNoRoomAvailable)
	}

	if err := h.WakeUpPatient(inBed); err != nil {
		t.Fatalf("WakeUpPatient error: %v", err)
	}
	if inBed.State() != Awake || inBed.Room() != nil {
		t.Errorf("woken patient = %v in room %v, want %v without room", inBed.State(), inBed.Room(), Awake)
	}
	if room.State() != Available || len(room.Occupants()) != 0 {
		t.Errorf("room = %v with %d occupants, want %v and empty", room.State(), len(room.Occupants()), Available)
	}

	got, err := h.AssignRoom(inHallway)
	if err != nil {
		t.Fatalf("AssignRoom error: %v", err)
	}
	if got != room || inHallway.State() != AsleepInBed {
		t.Errorf("AssignRoom = room %v, patient %v, want room %d and %v", got, inHallway.State(), room.Number(), AsleepInBed)
	}
	// Getting a room later does not create a new episode.
	if n := len(h.History()); n != 2 {
		t.Errorf("len(History()) = %d, want 2", n)
	}
}

func TestHospitalHandleSleepAttackErrors(t *testing.T) {
	h := NewHospital("Test Hospital") // no rooms and no staff
	p := newTestPatient(t, "P1")

	_, nilErr := h.HandleSleepAttack(nil, "cafeteria")
	_, noStaffErr := h.HandleSleepAttack(p, "cafeteria")
	_, asleepErr := h.HandleSleepAttack(p, "cafeteria")

	tests := []struct {
		name string
		err  error
		want error
	}{
		{"nil patient", nilErr, ErrInvalidData},
		{"no staff", noStaffErr, ErrNoStaff},
		{"already asleep", asleepErr, ErrAlreadyAsleep},
		{"WakeUpPatient(nil)", h.WakeUpPatient(nil), ErrInvalidData},
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

// admitTestPatient builds a patient with the given level, admits it to h and
// returns it, or stops the test.
func admitTestPatient(t *testing.T, h *Hospital, id string, level NarcolepsyLevel) *Patient {
	t.Helper()
	p, err := NewPatient(id, "Patient "+id, 30, level)
	if err != nil {
		t.Fatalf("NewPatient(%q) error: %v", id, err)
	}
	if err := h.AdmitPatient(p); err != nil {
		t.Fatalf("AdmitPatient(%s) error: %v", id, err)
	}
	return p
}

func TestPatientsInHallway(t *testing.T) {
	h, patients := newTestHospital(t, 1)
	inBed, inHallway := patients[0], patients[1] // patients[2] stays awake

	if got := h.PatientsInHallway(); got == nil || len(got) != 0 {
		t.Errorf("PatientsInHallway() before any attack = %v, want an empty non-nil slice", got)
	}

	if _, err := h.HandleSleepAttack(inBed, "cafeteria"); err != nil {
		t.Fatalf("HandleSleepAttack(P1) error: %v", err)
	}
	if _, err := h.HandleSleepAttack(inHallway, "library"); !errors.Is(err, ErrNoRoomAvailable) {
		t.Fatalf("HandleSleepAttack(P2) error = %v, want %v", err, ErrNoRoomAvailable)
	}

	got := h.PatientsInHallway()
	if len(got) != 1 || got[0] != inHallway {
		t.Errorf("PatientsInHallway() = %v, want only %s", got, inHallway.ID())
	}

	// Once awake, the patient is no longer in the hallway.
	if err := h.WakeUpPatient(inHallway); err != nil {
		t.Fatalf("WakeUpPatient(%s) error: %v", inHallway.ID(), err)
	}
	if got := h.PatientsInHallway(); len(got) != 0 {
		t.Errorf("PatientsInHallway() after waking up = %v, want empty", got)
	}
}

// runSleepAttacks makes n sleep attacks through h, cycling over patients and
// waking each patient up right after, or stops the test. h must have staff
// and at least one room, so every attack gets a bed.
func runSleepAttacks(t *testing.T, h *Hospital, patients []*Patient, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		p := patients[i%len(patients)]
		if _, err := h.HandleSleepAttack(p, "cafeteria"); err != nil {
			t.Fatalf("attack #%d: HandleSleepAttack(%s) error: %v", i+1, p.ID(), err)
		}
		if err := h.WakeUpPatient(p); err != nil {
			t.Fatalf("attack #%d: WakeUpPatient(%s) error: %v", i+1, p.ID(), err)
		}
	}
}

func TestSevereReport(t *testing.T) {
	type patientData struct {
		id    string
		level NarcolepsyLevel
	}
	tests := []struct {
		name     string
		patients []patientData
		attacks  []string       // IDs of the patients that fall asleep, in order
		want     map[string]int // patient ID -> episodes today
	}{
		{
			name:     "counts the episodes of each severe patient",
			patients: []patientData{{"P1", Severe}, {"P2", Severe}},
			attacks:  []string{"P1", "P2", "P1"},
			want:     map[string]int{"P1": 2, "P2": 1},
		},
		{
			name:     "severe patient without episodes appears with 0",
			patients: []patientData{{"P1", Severe}, {"P2", Severe}},
			attacks:  []string{"P1"},
			want:     map[string]int{"P1": 1, "P2": 0},
		},
		{
			name:     "mild and moderate patients are excluded",
			patients: []patientData{{"P1", Mild}, {"P2", Moderate}, {"P3", Severe}},
			attacks:  []string{"P1", "P2", "P2", "P3"},
			want:     map[string]int{"P3": 1},
		},
		{
			name:     "no severe patients gives an empty report",
			patients: []patientData{{"P1", Mild}},
			attacks:  []string{"P1"},
			want:     map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHospitalWithStaff(t, newTestRoom(t, 101, 1))
			byID := make(map[string]*Patient)
			for _, pd := range tt.patients {
				byID[pd.id] = admitTestPatient(t, h, pd.id, pd.level)
			}
			var attacked []*Patient
			for _, id := range tt.attacks {
				attacked = append(attacked, byID[id])
			}
			runSleepAttacks(t, h, attacked, len(attacked))

			// The report uses *Patient keys; comparing by ID gives readable
			// messages instead of memory addresses.
			got := make(map[string]int)
			for p, n := range h.SevereReport() {
				got[p.ID()] = n
			}
			if !maps.Equal(got, tt.want) {
				t.Errorf("SevereReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSevereReport_OnlyToday(t *testing.T) {
	h := newTestHospitalWithStaff(t, newTestRoom(t, 101, 1))
	p := admitTestPatient(t, h, "P1", Severe)
	runSleepAttacks(t, h, []*Patient{p}, 1)

	// An episode from yesterday must not be counted.
	yesterday := time.Now().AddDate(0, 0, -1)
	old := newEpisodeRecord("OLD-E1", yesterday, p, newTestDoctor(t, "D-99"), "library", nil)
	if err := h.RegisterEpisode(old); err != nil {
		t.Fatalf("RegisterEpisode error: %v", err)
	}

	if n := h.SevereReport()[p]; n != 1 {
		t.Errorf("SevereReport()[%s] = %d, want 1", p.ID(), n)
	}
}

func TestMyEpisodes(t *testing.T) {
	// The staff is D-01, D-02, O-01 and attack i hits patient Pi, so the
	// round-robin decides who attends whom.
	tests := []struct {
		name    string
		attacks int
		want    map[string][]string // doctor ID -> patients they attended
	}{
		{"one attack", 1, map[string][]string{"D-01": {"P1"}, "D-02": {}}},
		{"one turn of the whole staff", 3, map[string][]string{"D-01": {"P1"}, "D-02": {"P2"}}},
		{"second turn", 5, map[string][]string{"D-01": {"P1", "P4"}, "D-02": {"P2", "P5"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHospitalWithRooms(t, newTestRoom(t, 101, 1))
			doctors := []*Doctor{newTestDoctor(t, "D-01"), newTestDoctor(t, "D-02")}
			for _, d := range doctors {
				if err := h.HireDoctor(d); err != nil {
					t.Fatalf("HireDoctor(%s) error: %v", d.ID(), err)
				}
			}
			if err := h.HireStaff(newTestOrderly(t, "O-01")); err != nil {
				t.Fatalf("HireStaff error: %v", err)
			}
			var patients []*Patient
			for i := 1; i <= tt.attacks; i++ {
				patients = append(patients, admitTestPatient(t, h, fmt.Sprintf("P%d", i), Moderate))
			}
			runSleepAttacks(t, h, patients, tt.attacks)

			for _, d := range doctors {
				var got []string
				for _, e := range d.MyEpisodes() {
					if e.AttendedBy().ID() != d.ID() {
						t.Errorf("%s has episode %s attended by %s", d.ID(), e.ID(), e.AttendedBy().ID())
					}
					got = append(got, e.Patient().ID())
				}
				if want := tt.want[d.ID()]; !slices.Equal(got, want) {
					t.Errorf("%s MyEpisodes() patients = %v, want %v", d.ID(), got, want)
				}
			}
		})
	}
}

func TestHistoryAttendedByDoctorAndOrderly(t *testing.T) {
	// The staff is D-01 (*Doctor) and O-01 (*Orderly), taking turns.
	tests := []struct {
		name    string
		attacks int
		want    map[string]int // staff ID -> records in History()
	}{
		{"one attack each", 2, map[string]int{"D-01": 1, "O-01": 1}},
		{"two attacks each", 4, map[string]int{"D-01": 2, "O-01": 2}},
		{"odd number of attacks", 5, map[string]int{"D-01": 3, "O-01": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, patients := newTestHospital(t, 1)
			runSleepAttacks(t, h, patients, tt.attacks)

			// AttendedBy() is an Attender: ID() works for both types, so no
			// type switch is needed.
			got := make(map[string]int)
			for _, rec := range h.History() {
				got[rec.AttendedBy().ID()]++
			}
			if !maps.Equal(got, tt.want) {
				t.Errorf("records per staff member = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorCases(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T) error
		want error
	}{
		{
			name: "WakeUp on awake patient",
			run: func(t *testing.T) error {
				return newTestPatient(t, "P1").WakeUp()
			},
			want: ErrAlreadyAwake,
		},
		{
			name: "SufferSleepAttack on patient asleep in hallway",
			run: func(t *testing.T) error {
				return newAsleepPatient(t, "P1").SufferSleepAttack("library")
			},
			want: ErrAlreadyAsleep,
		},
		{
			name: "SufferSleepAttack on patient asleep in bed",
			run: func(t *testing.T) error {
				h, patients := newTestHospital(t, 1)
				p := patients[0]
				if _, err := h.HandleSleepAttack(p, "cafeteria"); err != nil {
					t.Fatalf("HandleSleepAttack error: %v", err)
				}
				return p.SufferSleepAttack("library")
			},
			want: ErrAlreadyAsleep,
		},
		{
			name: "AdmitPatient with duplicate ID",
			run: func(t *testing.T) error {
				h, _ := newTestHospital(t, 0) // already has P1
				return h.AdmitPatient(newTestPatient(t, "P1"))
			},
			want: ErrDuplicateID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(t); !errors.Is(err, tt.want) {
				t.Errorf("error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSameDay(t *testing.T) {
	base := time.Date(2026, time.October, 20, 0, 0, 0, 0, time.Local)
	tests := []struct {
		name string
		b    time.Time
		want bool
	}{
		{"same moment", base, true},
		{"end of the same day", base.Add(23*time.Hour + 59*time.Minute), true},
		{"next day", base.AddDate(0, 0, 1), false},
		{"same day and month, other year", base.AddDate(-1, 0, 0), false},
	}
	for _, tt := range tests {
		if got := sameDay(base, tt.b); got != tt.want {
			t.Errorf("%s: sameDay = %t, want %t", tt.name, got, tt.want)
		}
	}
}

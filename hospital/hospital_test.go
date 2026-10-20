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

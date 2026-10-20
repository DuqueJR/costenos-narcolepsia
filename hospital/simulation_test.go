package hospital

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestSimulate(t *testing.T) {
	const numPatients, rounds = 10, 3
	// Fewer rooms than patients, so some attacks find no room.
	h := newTestHospitalWithStaff(t, newTestRoom(t, 101, 1), newTestRoom(t, 102, 2))
	for i := 1; i <= numPatients; i++ {
		admitTestPatient(t, h, fmt.Sprintf("P%02d", i), Severe)
	}

	res, err := h.Simulate(SimConfig{
		Rounds:    rounds,
		MaxDelay:  2 * time.Millisecond,
		Locations: []string{"cafeteria", "parking lot", "beach"},
	})
	if err != nil {
		t.Fatalf("Simulate error: %v", err)
	}

	// Simulate already waited for every goroutine, so reading is safe here.
	want := numPatients * rounds
	if n := len(h.History()); n != want {
		t.Errorf("len(History()) = %d, want %d (patients x rounds)", n, want)
	}
	if res.Attacks != want {
		t.Errorf("Attacks = %d, want %d", res.Attacks, want)
	}
	if res.RoomFailures > res.Attacks {
		t.Errorf("RoomFailures = %d, more than Attacks = %d", res.RoomFailures, res.Attacks)
	}
	if res.Errors != 0 {
		t.Errorf("Errors = %d, want 0", res.Errors)
	}
	// Every round ends with a wake up, so everybody is awake and every
	// room is empty at the end.
	for _, p := range h.Patients() {
		if p.State() != Awake {
			t.Errorf("patient %s state = %v, want %v", p.ID(), p.State(), Awake)
		}
	}
	for _, r := range h.Rooms() {
		if n := len(r.Occupants()); n != 0 {
			t.Errorf("room %d has %d occupants, want 0", r.Number(), n)
		}
	}
}

func TestSimulateInvalidConfig(t *testing.T) {
	valid := SimConfig{Rounds: 1, MaxDelay: time.Millisecond, Locations: []string{"cafeteria"}}
	tests := []struct {
		name   string
		change func(cfg *SimConfig)
	}{
		{"zero rounds", func(cfg *SimConfig) { cfg.Rounds = 0 }},
		{"negative max delay", func(cfg *SimConfig) { cfg.MaxDelay = -time.Millisecond }},
		{"no locations", func(cfg *SimConfig) { cfg.Locations = nil }},
		{"blank location", func(cfg *SimConfig) { cfg.Locations = []string{"cafeteria", "  "} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := newTestHospital(t, 1)
			cfg := valid
			tt.change(&cfg)
			if _, err := h.Simulate(cfg); !errors.Is(err, ErrInvalidData) {
				t.Errorf("error = %v, want %v", err, ErrInvalidData)
			}
			if n := len(h.History()); n != 0 {
				t.Errorf("len(History()) = %d, want 0: an invalid config must not start the simulation", n)
			}
		})
	}
}

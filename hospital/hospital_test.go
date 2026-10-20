package hospital

import (
	"fmt"
	"testing"
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

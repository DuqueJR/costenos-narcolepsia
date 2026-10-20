package hospital

import "fmt"

// Patient is a person with narcolepsy treated at the hospital. It embeds
// Person by value, so ID, Name and Age are promoted to Patient. Its state
// only changes through methods, never by touching the fields directly.
type Patient struct {
	Person
	level           NarcolepsyLevel
	state           PatientState
	currentLocation string
	room            *Room
	assignedDoctor  *Doctor
}

// NewPatient validates the data and builds an awake Patient. It returns an
// error wrapping ErrInvalidData if the person data or the level are invalid.
func NewPatient(id, name string, age int, level NarcolepsyLevel) (*Patient, error) {
	person, err := newPerson(id, name, age)
	if err != nil {
		return nil, fmt.Errorf("new patient: %w", err)
	}
	if level < Mild || level > Severe {
		return nil, fmt.Errorf("new patient %s: invalid narcolepsy level %d: %w", person.ID(), int(level), ErrInvalidData)
	}
	return &Patient{Person: person, level: level, state: Awake}, nil
}

// SufferSleepAttack makes the patient fall asleep at location, in the
// hallway and without a room. It returns an error wrapping ErrAlreadyAsleep
// if the patient is already asleep.
func (p *Patient) SufferSleepAttack(location string) error {
	if p.state != Awake {
		return fmt.Errorf("patient %s: %w", p.ID(), ErrAlreadyAsleep)
	}
	p.state = AsleepInHallway
	p.currentLocation = location
	p.room = nil
	return nil
}

// WakeUp makes the patient awake again and leaves them without a room. It
// returns an error wrapping ErrAlreadyAwake if the patient is already awake.
func (p *Patient) WakeUp() error {
	if p.state == Awake {
		return fmt.Errorf("patient %s: %w", p.ID(), ErrAlreadyAwake)
	}
	p.state = Awake
	p.room = nil
	return nil
}

// State returns whether the patient is awake or where they are asleep.
func (p *Patient) State() PatientState {
	return p.state
}

// Level returns how severe the narcolepsy of the patient is.
func (p *Patient) Level() NarcolepsyLevel {
	return p.level
}

// Location returns the place where the patient had the last sleep attack.
func (p *Patient) Location() string {
	return p.currentLocation
}

// Room returns the room where the patient sleeps, or nil if they have none.
func (p *Patient) Room() *Room {
	return p.room
}

// AssignedDoctor returns the doctor in charge of the patient, or nil if no
// doctor has diagnosed them yet.
func (p *Patient) AssignedDoctor() *Doctor {
	return p.assignedDoctor
}

// assignDoctor makes d the doctor in charge of the patient. Only Doctor
// calls it from DiagnosePatient, so the doctor and the patient always agree
// on who is in charge.
func (p *Patient) assignDoctor(d *Doctor) {
	p.assignedDoctor = d
}

// moveToBed puts the asleep patient in a bed of room r. Only Room calls it,
// so the patient and the room always agree on who is where.
func (p *Patient) moveToBed(r *Room) {
	p.state = AsleepInBed
	p.room = r
}

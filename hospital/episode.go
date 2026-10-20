package hospital

import (
	"fmt"
	"time"
)

// EpisodeRecord is the record of one sleep attack: who fell asleep, where,
// when, who attended them and in which room they ended up. A nil room means
// the patient stayed asleep in the hallway. Its fields are unexported and it
// is only built by an Attender, so a record cannot change once created.
type EpisodeRecord struct {
	id         string
	at         time.Time
	patient    *Patient
	attendedBy Attender
	location   string
	room       *Room
}

// newEpisodeRecord builds an EpisodeRecord. It is unexported because only
// the Attender implementations of this package create records.
func newEpisodeRecord(id string, at time.Time, p *Patient, by Attender, location string, room *Room) EpisodeRecord {
	return EpisodeRecord{
		id:         id,
		at:         at,
		patient:    p,
		attendedBy: by,
		location:   location,
		room:       room,
	}
}

// ID returns the identifier of the record.
func (e EpisodeRecord) ID() string {
	return e.id
}

// At returns when the episode was attended.
func (e EpisodeRecord) At() time.Time {
	return e.at
}

// Patient returns the patient who had the sleep attack.
func (e EpisodeRecord) Patient() *Patient {
	return e.patient
}

// AttendedBy returns the staff member who attended the episode.
func (e EpisodeRecord) AttendedBy() Attender {
	return e.attendedBy
}

// Location returns the place where the patient fell asleep.
func (e EpisodeRecord) Location() string {
	return e.location
}

// Room returns the room the patient was taken to, or nil if they stayed in
// the hallway.
func (e EpisodeRecord) Room() *Room {
	return e.room
}

// Summary returns the record in one line, for example
// "2026-10-19 14:15 P-001 cafeteria -> room 101 (Dr. Karen Ospina)",
// or "-> hallway" instead of the room when the patient has none.
func (e EpisodeRecord) Summary() string {
	destination := "hallway"
	if e.room != nil {
		destination = fmt.Sprintf("room %d", e.room.Number())
	}
	return fmt.Sprintf("%s %s %s -> %s (%s)",
		e.at.Format("2006-01-02 15:04"), e.patient.ID(), e.location, destination, e.attendedBy.Name())
}

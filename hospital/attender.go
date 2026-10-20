package hospital

// Attender is anyone on the staff who can attend a patient during a sleep
// attack. Doctor and Orderly implement it, and the hospital keeps its staff
// in a []Attender and dispatches every emergency through this interface,
// without knowing the concrete type behind it.
type Attender interface {
	// ID returns the identifier of the staff member.
	ID() string
	// Name returns the full name of the staff member.
	Name() string
	// Attend handles patient p, asleep at location, and returns the record
	// of the episode. It returns an error if p cannot be attended.
	Attend(p *Patient, location string) (EpisodeRecord, error)
}

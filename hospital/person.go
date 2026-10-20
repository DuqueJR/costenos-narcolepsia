package hospital

import (
	"fmt"
	"strings"
)

// Person is the identity shared by patients and staff. Its fields are
// unexported, so other packages can only read them through the getters.
// A Person is never built directly: it is embedded by value in Patient,
// Doctor and Orderly, whose constructors call newPerson.
type Person struct {
	id   string
	name string
	age  int
}

// newPerson validates the data and builds a Person. It returns an error
// wrapping ErrInvalidData if id or name are blank or age is not positive.
func newPerson(id, name string, age int) (Person, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if id == "" {
		return Person{}, fmt.Errorf("new person: empty id: %w", ErrInvalidData)
	}
	if name == "" {
		return Person{}, fmt.Errorf("new person %s: empty name: %w", id, ErrInvalidData)
	}
	if age <= 0 {
		return Person{}, fmt.Errorf("new person %s: age %d must be positive: %w", id, age, ErrInvalidData)
	}
	return Person{id: id, name: name, age: age}, nil
}

// ID returns the identifier of the person.
func (p Person) ID() string {
	return p.id
}

// Name returns the full name of the person.
func (p Person) Name() string {
	return p.name
}

// Age returns the age of the person in years.
func (p Person) Age() int {
	return p.age
}

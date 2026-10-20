package main

import "costenos-narcolepsia/hospital"

// newSeedHospital builds the initial state of the scenario: 2 doctors,
// 1 orderly, 3 rooms with one bed each and 5 awake patients, with the same
// IDs and names used by main.go in the root module. No sleep attack happens
// here; the user triggers them from the GUI. It returns the first error found.
//
// This duplicates data, not logic: the GUI is another program in another
// module, so it cannot call buildHospital from the root main package. Every
// rule still lives in hospital/, and both programs only call its methods.
func newSeedHospital() (*hospital.Hospital, error) {
	h := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")

	karen, err := hospital.NewDoctor("D-01", "Karen Ospina", 41, "Neurology")
	if err != nil {
		return nil, err
	}
	if err := h.HireDoctor(karen); err != nil {
		return nil, err
	}

	luis, err := hospital.NewDoctor("D-02", "Luis Carlos Barrios", 52, "Sleep Medicine")
	if err != nil {
		return nil, err
	}
	if err := h.HireDoctor(luis); err != nil {
		return nil, err
	}

	jhon, err := hospital.NewOrderly("O-01", "Jhon Jairo Rada", 29)
	if err != nil {
		return nil, err
	}
	if err := h.HireStaff(jhon); err != nil {
		return nil, err
	}

	for _, number := range []int{101, 102, 103} {
		r, err := hospital.NewRoom(number, 1)
		if err != nil {
			return nil, err
		}
		if err := h.AddRoom(r); err != nil {
			return nil, err
		}
	}

	for _, data := range []struct {
		id    string
		name  string
		age   int
		level hospital.NarcolepsyLevel
	}{
		{"P-001", "Carlos Polo", 34, hospital.Moderate},
		{"P-002", "Yeimy Padilla", 27, hospital.Severe},
		{"P-003", "Dayana Mercado", 45, hospital.Mild},
		{"P-004", "Wilfrido Berrio", 63, hospital.Severe},
		{"P-005", "Kevin Arrieta", 19, hospital.Moderate},
	} {
		p, err := hospital.NewPatient(data.id, data.name, data.age, data.level)
		if err != nil {
			return nil, err
		}
		if err := h.AdmitPatient(p); err != nil {
			return nil, err
		}
	}
	return h, nil
}

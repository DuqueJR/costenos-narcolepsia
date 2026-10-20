// Command costenos-narcolepsia builds the hospital scenario and prints its results.
package main

import (
	"errors"
	"fmt"
	"os"

	"costenos-narcolepsia/hospital"
)

func main() {
	// 1. Build the hospital: staff, rooms and patients.
	h := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")

	karen, err := hospital.NewDoctor("D-01", "Karen Ospina", 41, "Neurology")
	exitOnError(err)
	exitOnError(h.HireDoctor(karen))

	luis, err := hospital.NewDoctor("D-02", "Luis Carlos Barrios", 52, "Sleep Medicine")
	exitOnError(err)
	exitOnError(h.HireDoctor(luis))

	jhon, err := hospital.NewOrderly("O-01", "Jhon Jairo Rada", 29)
	exitOnError(err)
	exitOnError(h.HireStaff(jhon))

	for _, number := range []int{101, 102, 103} {
		r, err := hospital.NewRoom(number, 1)
		exitOnError(err)
		exitOnError(h.AddRoom(r))
	}

	var patients []*hospital.Patient
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
		exitOnError(err)
		exitOnError(h.AdmitPatient(p))
		patients = append(patients, p)
	}

	fmt.Printf("== %s ==\n", h.Name())
	fmt.Printf("  %d doctors, %d staff in total, %d rooms, %d patients admitted\n",
		len(h.Doctors()), len(h.Staff()), len(h.Rooms()), len(h.Patients()))

	// 2. Four patients suffer a sleep attack. The fourth one finds no room.
	fmt.Println()
	fmt.Println("== Sleep attacks ==")
	attacks := []struct {
		patient  *hospital.Patient
		location string
	}{
		{patients[0], "cafeteria"},
		{patients[1], "radiology waiting room"},
		{patients[2], "parking lot"},
		{patients[3], "hallway 2, second floor"},
	}
	var noRoomErr error // the "no room available" error, kept for the final report
	for _, a := range attacks {
		rec, err := h.HandleSleepAttack(a.patient, a.location)
		switch {
		case err == nil:
			fmt.Println("  " + rec.Summary())
		case errors.Is(err, hospital.ErrNoRoomAvailable):
			// The episode was registered anyway: the patient sleeps in the hallway.
			fmt.Println("  " + rec.Summary())
			fmt.Printf("  ! %s %s could not get a bed and stays in the hallway (%v)\n",
				a.patient.ID(), a.patient.Name(), err)
			noRoomErr = err
		default:
			exitOnError(err)
		}
	}

	// 3. Right after the failure, the patient is visible in the hallway.
	fmt.Println()
	printHallway(h)

	// 4. P-002 wakes up and frees a bed, which then goes to P-004.
	fmt.Println()
	fmt.Println("== P-002 wakes up ==")
	sleeper, hallwayPatient := patients[1], patients[3]
	freed := sleeper.Room() // read before waking up: WakeUpPatient leaves the patient without a room
	exitOnError(h.WakeUpPatient(sleeper))
	if freed != nil {
		fmt.Printf("  %s %s woke up and released room %d\n", sleeper.ID(), sleeper.Name(), freed.Number())
	} else {
		fmt.Printf("  %s %s woke up (had no room)\n", sleeper.ID(), sleeper.Name())
	}
	reassigned, err := h.AssignRoom(hallwayPatient)
	exitOnError(err)
	fmt.Printf("  %s %s -> room %d\n", hallwayPatient.ID(), hallwayPatient.Name(), reassigned.Number())

	// 5. The four required queries.
	fmt.Println()
	printHallway(h)

	for _, d := range h.Doctors() {
		fmt.Println()
		fmt.Printf("== Episodes attended by Dr. %s (%s) ==\n", d.Name(), d.Specialty())
		episodes := d.MyEpisodes()
		if len(episodes) == 0 {
			fmt.Println("  (none)")
		}
		for _, rec := range episodes {
			fmt.Println("  " + rec.Summary())
		}
	}

	fmt.Println()
	fmt.Println("== Room assignment ==")
	fmt.Printf("  first attempt for %s: %v\n", hallwayPatient.ID(), noRoomErr)
	fmt.Printf("  after %s woke up:    %s -> room %d\n", sleeper.ID(), hallwayPatient.ID(), reassigned.Number())

	fmt.Println()
	fmt.Println("== Severity report (Severe patients) ==")
	for p, episodes := range h.SevereReport() {
		fmt.Printf("  %s %s  %d episodes today\n", p.ID(), p.Name(), episodes)
	}
}

// printHallway prints the patients asleep in a hallway, or "(none)".
func printHallway(h *hospital.Hospital) {
	fmt.Println("== Patients asleep in hallway ==")
	hallway := h.PatientsInHallway()
	if len(hallway) == 0 {
		fmt.Println("  (none)")
	}
	for _, p := range hallway {
		fmt.Printf("  %s %s (%s)\n", p.ID(), p.Name(), p.Location())
	}
}

// exitOnError prints err and ends the program with status 1 when err is not
// nil. It is only used for errors the scenario does not expect.
func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "unexpected error:", err)
		os.Exit(1)
	}
}

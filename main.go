// Command costenos-narcolepsia builds the hospital scenario and prints its results.
package main

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"

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
	printHeader("Sleep attacks")
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
			fmt.Printf("  ! %s could not get a bed and stays in the hallway (%v)\n",
				patientColumns(a.patient), err)
			noRoomErr = err
		default:
			exitOnError(err)
		}
	}

	// 3. Right after the failure, the patient is visible in the hallway.
	printHallway(h)

	// 4. P-002 wakes up and frees a bed, which then goes to P-004.
	printHeader("P-002 wakes up")
	sleeper, hallwayPatient := patients[1], patients[3]
	freed := sleeper.Room() // read before waking up: WakeUpPatient leaves the patient without a room
	exitOnError(h.WakeUpPatient(sleeper))
	if freed != nil {
		fmt.Printf("  %s woke up and released room %d\n", patientColumns(sleeper), freed.Number())
	} else {
		fmt.Printf("  %s woke up (had no room)\n", patientColumns(sleeper))
	}
	reassigned, err := h.AssignRoom(hallwayPatient)
	exitOnError(err)
	fmt.Printf("  %s -> room %d\n", patientColumns(hallwayPatient), reassigned.Number())

	// 5. The four required queries.
	printHallway(h)
	for _, d := range h.Doctors() {
		printDoctorEpisodes(d)
	}
	printRoomAssignment(hallwayPatient, noRoomErr, reassigned)
	printSevereReport(h)
}

// printHeader prints a blank line and the title of a block of the output.
func printHeader(title string) {
	fmt.Printf("\n== %s ==\n", title)
}

// patientColumns returns the ID and the name of p padded to a fixed width,
// so every line that starts with a patient keeps its columns aligned.
func patientColumns(p *hospital.Patient) string {
	return fmt.Sprintf("%-6s %-16s", p.ID(), p.Name())
}

// printHallway prints the patients asleep in a hallway, or "(none)".
func printHallway(h *hospital.Hospital) {
	printHeader("Patients asleep in hallway")
	hallway := h.PatientsInHallway()
	if len(hallway) == 0 {
		fmt.Println("  (none)")
	}
	for _, p := range hallway {
		fmt.Printf("  %s (%s)\n", patientColumns(p), p.Location())
	}
}

// printDoctorEpisodes prints every episode attended by doctor d, or "(none)".
func printDoctorEpisodes(d *hospital.Doctor) {
	printHeader(fmt.Sprintf("Episodes attended by Dr. %s (%s)", d.Name(), d.Specialty()))
	episodes := d.MyEpisodes()
	if len(episodes) == 0 {
		fmt.Println("  (none)")
	}
	for _, rec := range episodes {
		fmt.Println("  " + rec.Summary())
	}
}

// printRoomAssignment prints how the room assignment of patient p went: the
// error of the first attempt and the room it got on the second one.
func printRoomAssignment(p *hospital.Patient, firstErr error, room *hospital.Room) {
	printHeader("Room assignment")
	fmt.Printf("  %s first attempt:  %v\n", patientColumns(p), firstErr)
	fmt.Printf("  %s second attempt: room %d\n", patientColumns(p), room.Number())
}

// printSevereReport prints the result of SevereReport sorted by patient ID.
// A map has no order, so its keys are sorted first to get the same output on
// every run.
func printSevereReport(h *hospital.Hospital) {
	printHeader("Severity report (Severe patients)")
	report := h.SevereReport()
	severe := make([]*hospital.Patient, 0, len(report))
	for p := range report {
		severe = append(severe, p)
	}
	slices.SortFunc(severe, func(a, b *hospital.Patient) int {
		return cmp.Compare(a.ID(), b.ID())
	})
	for _, p := range severe {
		fmt.Printf("  %s %d episode(s) today\n", patientColumns(p), report[p])
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

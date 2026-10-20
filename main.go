// Command costenos-narcolepsia builds the hospital scenario and prints its results.
package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"time"

	"costenos-narcolepsia/hospital"
)

func main() {
	sim := flag.Bool("sim", false, "run the concurrent simulation instead of the scenario")
	flag.Parse()

	// 1. Build the hospital: staff, rooms and patients.
	h, patients, err := buildHospital()
	exitOnError(err)

	fmt.Printf("== %s ==\n", h.Name())
	fmt.Printf("  %d doctors, %d staff in total, %d rooms, %d patients admitted\n",
		len(h.Doctors()), len(h.Staff()), len(h.Rooms()), len(h.Patients()))

	if *sim {
		runSimulation(h)
		return
	}
	runScenario(h, patients)
}

// buildHospital builds the hospital of the scenario: 2 doctors, 1 orderly,
// 3 rooms with one bed each and 5 awake patients. It returns the patients in
// the order they were admitted, or the first error found.
func buildHospital() (*hospital.Hospital, []*hospital.Patient, error) {
	h := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")

	karen, err := hospital.NewDoctor("D-01", "Karen Ospina", 41, "Neurology")
	if err != nil {
		return nil, nil, err
	}
	if err := h.HireDoctor(karen); err != nil {
		return nil, nil, err
	}

	luis, err := hospital.NewDoctor("D-02", "Luis Carlos Barrios", 52, "Sleep Medicine")
	if err != nil {
		return nil, nil, err
	}
	if err := h.HireDoctor(luis); err != nil {
		return nil, nil, err
	}

	jhon, err := hospital.NewOrderly("O-01", "Jhon Jairo Rada", 29)
	if err != nil {
		return nil, nil, err
	}
	if err := h.HireStaff(jhon); err != nil {
		return nil, nil, err
	}

	for _, number := range []int{101, 102, 103} {
		r, err := hospital.NewRoom(number, 1)
		if err != nil {
			return nil, nil, err
		}
		if err := h.AddRoom(r); err != nil {
			return nil, nil, err
		}
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
		if err != nil {
			return nil, nil, err
		}
		if err := h.AdmitPatient(p); err != nil {
			return nil, nil, err
		}
		patients = append(patients, p)
	}
	return h, patients, nil
}

// runSimulation runs the concurrent simulation on h and prints its result
// and the four queries. Simulate has already waited for every goroutine when
// it returns, so the queries read a hospital that nobody is changing.
func runSimulation(h *hospital.Hospital) {
	cfg := hospital.SimConfig{
		Rounds:    3,
		MaxDelay:  50 * time.Millisecond,
		Locations: []string{"cafeteria", "parking lot", "radiology waiting room", "hallway 2, second floor"},
	}
	printHeader("Concurrent simulation")
	fmt.Printf("  %d patients x %d rounds, random pauses up to %v\n",
		len(h.Patients()), cfg.Rounds, cfg.MaxDelay)
	res, err := h.Simulate(cfg)
	exitOnError(err)
	fmt.Printf("  attacks: %d, without room: %d, errors: %d\n",
		res.Attacks, res.RoomFailures, res.Errors)

	printHallway(h)
	for _, d := range h.Doctors() {
		printDoctorEpisodes(d)
	}
	// Every round ends with WakeUpPatient, so the hallway is normally empty
	// and there is nobody to assign; the AssignRoom failures during the
	// simulation are the "without room" count above.
	printHeader("Room assignment (patients in hallway)")
	hallway := h.PatientsInHallway()
	if len(hallway) == 0 {
		fmt.Println("  (none)")
	}
	for _, p := range hallway {
		if r, err := h.AssignRoom(p); err != nil {
			fmt.Printf("  %s %v\n", patientColumns(p), err)
		} else {
			fmt.Printf("  %s -> room %d\n", patientColumns(p), r.Number())
		}
	}
	printSevereReport(h)
}

// runScenario runs the required scenario on h, whose patients are the ones
// returned by buildHospital, and prints every step.
func runScenario(h *hospital.Hospital, patients []*hospital.Patient) {
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

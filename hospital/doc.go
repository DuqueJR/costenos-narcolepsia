// Package hospital models the Hospital de los Costeños con Narcolepsia:
// patients, rooms, doctors and orderlies, and the sleep attack flow that
// ties them together. It holds all the business logic and never prints;
// package main only builds a scenario and shows the results.
//
// # Model
//
// Person holds the identity (ID, Name, Age) shared by everyone. It is
// embedded by value in Patient, Doctor and Orderly, so its methods are
// promoted to them and no identity field is repeated.
//
//   - Patient has a NarcolepsyLevel and a PatientState: Awake,
//     AsleepInHallway or AsleepInBed. While asleep in bed it knows its Room.
//   - Room has a fixed number of beds and a RoomState: Available while it
//     has a free bed, Occupied when it is full.
//   - Doctor keeps the patients in its care and the episodes it attended.
//   - Orderly attends emergencies too, but never diagnoses a patient.
//   - EpisodeRecord is the immutable record of one sleep attack.
//   - Hospital registers patients, staff and rooms and runs the flow.
//
// Doctor and Orderly implement the Attender interface. The hospital keeps
// its staff in a []Attender and dispatches every emergency through it,
// taking turns (round-robin), without knowing the concrete type.
//
// # Sleep attack flow
//
// Hospital.HandleSleepAttack makes the patient fall asleep in the hallway,
// puts them in the first room with a free bed, sends the next staff member
// to attend them and adds the episode to the history. If no room has a free
// bed, the episode is still registered and the patient stays in the
// hallway: the returned error wraps ErrNoRoomAvailable. Hospital.WakeUpPatient
// releases the bed so another patient can take it.
//
// # Errors
//
// Fields are unexported and state only changes through methods. Methods
// that can fail return an error that wraps one of the sentinel errors in
// this package, such as ErrNoRoomAvailable, so callers check them with
// errors.Is.
//
// # Concurrency
//
// Every exported method of Hospital takes its mutex, so several goroutines
// can call them at the same time. The entities do not have a mutex of their
// own, so while goroutines are running the following contract applies:
//
//   - Patient, Room, Doctor and Orderly are only modified through the
//     methods of Hospital (HandleSleepAttack, AssignRoom, WakeUpPatient...),
//     never by calling their own methods directly.
//   - Doctor.MyEpisodes and the getters of the entities (Patient.State,
//     Room.Occupants...) are read only when no goroutine is writing, for
//     example after wg.Wait().
package hospital

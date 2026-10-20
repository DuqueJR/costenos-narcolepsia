package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// levels are the narcolepsy levels offered by the patient form, in order.
var levels = []hospital.NarcolepsyLevel{hospital.Mild, hospital.Moderate, hospital.Severe}

// specialties are the suggestions of the doctor form. The user can still
// type any other specialty.
var specialties = []string{"Neurology", "Sleep Medicine", "Psychiatry", "Pulmonology", "Internal Medicine"}

// showAdmitPatientForm asks for the data of a new patient and admits them
// with hospital.NewPatient and Hospital.AdmitPatient.
func (a *App) showAdmitPatientForm() {
	var taken []string
	for _, p := range a.h.Patients() {
		taken = append(taken, p.ID())
	}
	id := newRequiredEntry(nextID("P-", 3, taken))
	name := newRequiredEntry("")
	name.SetPlaceHolder("Nombre y apellido")
	age := newNumberEntry("")

	options := make([]string, 0, len(levels))
	for _, l := range levels {
		options = append(options, levelLabel(l))
	}
	level := widget.NewSelect(options, nil)
	level.SetSelectedIndex(0)

	items := []*widget.FormItem{
		widget.NewFormItem("ID", id),
		widget.NewFormItem("Nombre", name),
		widget.NewFormItem("Edad", age),
		widget.NewFormItem("Nivel", level),
	}
	a.showForm("Admitir paciente", items, func() error {
		p, err := hospital.NewPatient(text(id), text(name), number(age), levels[level.SelectedIndex()])
		if err != nil {
			return err
		}
		return a.h.AdmitPatient(p)
	})
}

// showHireDoctorForm asks for the data of a new doctor and hires them with
// hospital.NewDoctor and Hospital.HireDoctor.
func (a *App) showHireDoctorForm() {
	id := newRequiredEntry(nextID("D-", 2, a.staffIDs()))
	name := newRequiredEntry("")
	name.SetPlaceHolder("Nombre y apellido")
	age := newNumberEntry("")
	specialty := widget.NewSelectEntry(specialties)
	specialty.Validator = required
	specialty.SetText(specialties[0])

	items := []*widget.FormItem{
		widget.NewFormItem("ID", id),
		widget.NewFormItem("Nombre", name),
		widget.NewFormItem("Edad", age),
		widget.NewFormItem("Especialidad", specialty),
	}
	a.showForm("Contratar doctor", items, func() error {
		d, err := hospital.NewDoctor(text(id), text(name), number(age), strings.TrimSpace(specialty.Text))
		if err != nil {
			return err
		}
		return a.h.HireDoctor(d)
	})
}

// showHireOrderlyForm asks for the data of a new orderly and hires them with
// hospital.NewOrderly and Hospital.HireStaff.
func (a *App) showHireOrderlyForm() {
	id := newRequiredEntry(nextID("O-", 2, a.staffIDs()))
	name := newRequiredEntry("")
	name.SetPlaceHolder("Nombre y apellido")
	age := newNumberEntry("")

	items := []*widget.FormItem{
		widget.NewFormItem("ID", id),
		widget.NewFormItem("Nombre", name),
		widget.NewFormItem("Edad", age),
	}
	a.showForm("Contratar camillero", items, func() error {
		o, err := hospital.NewOrderly(text(id), text(name), number(age))
		if err != nil {
			return err
		}
		return a.h.HireStaff(o)
	})
}

// showAddRoomForm asks for the number and capacity of a new room and adds
// it with hospital.NewRoom and Hospital.AddRoom.
func (a *App) showAddRoomForm() {
	next := 1
	for _, r := range a.h.Rooms() {
		next = max(next, r.Number()+1)
	}
	roomNumber := newNumberEntry(strconv.Itoa(next))
	capacity := newNumberEntry("1")

	items := []*widget.FormItem{
		widget.NewFormItem("Número", roomNumber),
		widget.NewFormItem("Camas", capacity),
	}
	a.showForm("Agregar cuarto", items, func() error {
		r, err := hospital.NewRoom(number(roomNumber), number(capacity))
		if err != nil {
			return err
		}
		return a.h.AddRoom(r)
	})
}

// showForm shows a dialog with items. When the user confirms, it runs save
// and, if save fails, shows its error; otherwise every view is refreshed.
func (a *App) showForm(title string, items []*widget.FormItem, save func() error) {
	d := dialog.NewForm(title, "Guardar", "Cancelar", items, func(confirmed bool) {
		if !confirmed {
			return
		}
		if err := save(); err != nil {
			a.showError(fmt.Errorf("%s: %w", strings.ToLower(title), err))
			return
		}
		a.refreshAll()
	}, a.win)
	d.Resize(fyne.NewSize(440, 0))
	d.Show()
}

// staffIDs returns the IDs of every staff member, doctors included, because
// the hospital does not allow two staff members with the same ID.
func (a *App) staffIDs() []string {
	var ids []string
	for _, s := range a.h.Staff() {
		ids = append(ids, s.ID())
	}
	return ids
}

// nextID returns the first ID made of prefix and a number of digits digits
// that is not in taken, starting after the IDs of taken that already use
// prefix. For example nextID("P-", 3, ids) returns "P-006" when there are 5
// patients, and nextID("O-", 2, staff) ignores the IDs of the doctors.
func nextID(prefix string, digits int, taken []string) string {
	n := 1
	for _, id := range taken {
		if strings.HasPrefix(id, prefix) {
			n++
		}
	}
	for ; ; n++ {
		id := fmt.Sprintf("%s%0*d", prefix, digits, n)
		if !slices.Contains(taken, id) {
			return id
		}
	}
}

// newRequiredEntry returns an entry with value that cannot be left empty.
func newRequiredEntry(value string) *widget.Entry {
	e := widget.NewEntry()
	e.SetText(value)
	e.Validator = required
	return e
}

// newNumberEntry returns an entry with value that only accepts a positive
// whole number.
func newNumberEntry(value string) *widget.Entry {
	e := widget.NewEntry()
	e.SetText(value)
	e.Validator = positiveNumber
	return e
}

// required is an entry validator that rejects empty text. The form keeps its
// "Guardar" button disabled while a validator fails.
func required(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("obligatorio")
	}
	return nil
}

// positiveNumber is an entry validator that only accepts a whole number
// greater than zero. The domain checks the same rules again, so this only
// gives earlier feedback.
func positiveNumber(s string) error {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n <= 0 {
		return errors.New("debe ser un número entero positivo")
	}
	return nil
}

// text returns the text of e without spaces around it.
func text(e *widget.Entry) string {
	return strings.TrimSpace(e.Text)
}

// number returns the text of e as a number. The validator of the entry
// already checked it, so a failure returns 0, which the domain rejects.
func number(e *widget.Entry) int {
	n, _ := strconv.Atoi(text(e))
	return n
}

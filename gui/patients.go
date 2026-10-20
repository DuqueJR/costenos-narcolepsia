package main

import (
	"errors"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// patientsView is the Patients tab: the list of patients on the left and the
// detail of the selected patient, with its actions, on the right.
type patientsView struct {
	a        *App
	patients []*hospital.Patient // what the list shows, reloaded on every refresh
	selected *hospital.Patient   // nil until the user picks a patient

	list *widget.List
	card *widget.Card

	id, age, level, state, location, room, doctor *widget.Label

	attackLocation *widget.Entry
	attackBtn      *widget.Button
	wakeBtn        *widget.Button
	assignBtn      *widget.Button
}

// newPatientsTab builds the Patients tab of app a and registers its refresher.
func newPatientsTab(a *App) fyne.CanvasObject {
	v := &patientsView{a: a}

	// The list asks how many rows there are, creates a few rows and reuses
	// them, filling each one with the patient at position i.
	v.list = widget.NewList(
		func() int {
			return len(v.patients)
		},
		func() fyne.CanvasObject {
			return newPatientRow()
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			// obj is always a row made by the function above.
			obj.(*patientRow).set(v.patients[i])
		},
	)
	v.list.OnSelected = func(i widget.ListItemID) {
		v.selected = v.patients[i]
		v.refreshDetail()
	}

	v.id = widget.NewLabel("")
	v.age = widget.NewLabel("")
	v.level = widget.NewLabel("")
	v.state = widget.NewLabel("")
	v.location = widget.NewLabel("")
	v.room = widget.NewLabel("")
	v.doctor = widget.NewLabel("")
	info := widget.NewForm(
		widget.NewFormItem("ID", v.id),
		widget.NewFormItem("Edad", v.age),
		widget.NewFormItem("Nivel", v.level),
		widget.NewFormItem("Estado", v.state),
		widget.NewFormItem("Ubicación (último ataque)", v.location),
		widget.NewFormItem("Cuarto", v.room),
		widget.NewFormItem("Doctor asignado", v.doctor),
	)

	v.attackLocation = widget.NewEntry()
	v.attackLocation.SetPlaceHolder("Ubicación, p. ej. cafeteria")
	v.attackBtn = widget.NewButtonWithIcon("Sleep attack", theme.WarningIcon(), v.sleepAttack)
	v.wakeBtn = widget.NewButtonWithIcon("Wake up", theme.VisibilityIcon(), v.wakeUp)
	v.assignBtn = widget.NewButtonWithIcon("Assign room", theme.HomeIcon(), v.assignRoom)
	actions := container.NewVBox(
		widget.NewForm(widget.NewFormItem("Ubicación", v.attackLocation)),
		container.NewGridWithColumns(3, v.attackBtn, v.wakeBtn, v.assignBtn),
	)

	v.card = widget.NewCard("", "", container.NewVBox(info, widget.NewSeparator(), actions))
	a.addRefresher(v.refresh)

	split := container.NewHSplit(v.list, container.NewVScroll(v.card))
	split.Offset = 0.4
	return split
}

// refresh reloads the list and the detail from the hospital. App.refreshAll
// calls it after every action.
func (v *patientsView) refresh() {
	v.patients = v.a.h.Patients()
	v.list.Refresh()
	v.refreshDetail()
}

// refreshDetail shows the selected patient in the card and enables only the
// actions that make sense for their current state.
func (v *patientsView) refreshDetail() {
	p := v.selected
	if p == nil {
		v.card.SetTitle("Selecciona un paciente")
		for _, l := range []*widget.Label{v.id, v.age, v.level, v.state, v.location, v.room, v.doctor} {
			l.SetText("-")
		}
		v.attackBtn.Disable()
		v.wakeBtn.Disable()
		v.assignBtn.Disable()
		return
	}

	v.card.SetTitle(p.Name())
	v.id.SetText(p.ID())
	v.age.SetText(strconv.Itoa(p.Age()))
	v.level.SetText(p.Level().String())
	v.state.SetText(p.State().String())
	v.location.SetText(locationText(p.Location()))
	v.room.SetText(roomText(p.Room()))
	v.doctor.SetText(doctorText(p.AssignedDoctor()))

	setEnabled(v.attackBtn, p.State() == hospital.Awake)
	setEnabled(v.wakeBtn, p.State() != hospital.Awake)
	setEnabled(v.assignBtn, p.State() == hospital.AsleepInHallway)
}

// sleepAttack runs a sleep attack of the selected patient at the location
// typed in the entry, and keeps the room it got as the last assignment.
func (v *patientsView) sleepAttack() {
	p := v.selected
	location := strings.TrimSpace(v.attackLocation.Text)
	if location == "" {
		v.a.showError(errors.New("escribe la ubicación donde se durmió el paciente"))
		return
	}

	_, err := v.a.h.HandleSleepAttack(p, location)
	switch {
	case err == nil:
		v.a.lastAssign = &assignResult{patient: p, room: p.Room()}
	case errors.Is(err, hospital.ErrNoRoomAvailable):
		// The episode is registered anyway: the patient sleeps in the hallway.
		v.a.lastAssign = &assignResult{patient: p, err: err}
		v.a.showWarning("Episodio registrado, pero no hay cuarto: el paciente queda en el pasillo")
	default:
		// Nothing was registered, so it is not an assignment result.
		v.a.showError(err)
	}
	v.a.refreshAll()
}

// wakeUp wakes the selected patient up, releasing their room if they had one.
func (v *patientsView) wakeUp() {
	if err := v.a.h.WakeUpPatient(v.selected); err != nil {
		v.a.showError(err)
	}
	v.a.refreshAll()
}

// assignRoom tries to give a room to the selected patient, asleep in the
// hallway, and keeps the outcome as the last assignment.
func (v *patientsView) assignRoom() {
	p := v.selected
	r, err := v.a.h.AssignRoom(p)
	v.a.lastAssign = &assignResult{patient: p, room: r, err: err}
	switch {
	case errors.Is(err, hospital.ErrNoRoomAvailable):
		v.a.showWarning("No hay cuartos libres: el paciente sigue en el pasillo")
	case err != nil:
		v.a.showError(err)
	}
	v.a.refreshAll()
}

// patientRow is one row of the patient list: an avatar with initials, the
// name and ID, and a pill colored by state. The list creates only a few rows
// and reuses them, so set changes what a row shows.
type patientRow struct {
	widget.BaseWidget
	avatar    *canvas.Circle
	initials  *canvas.Text
	name      *widget.Label
	id        *widget.Label
	pill      *canvas.Rectangle
	stateText *canvas.Text
}

// newPatientRow returns an empty row, ready to be filled with set.
func newPatientRow() *patientRow {
	r := &patientRow{
		avatar:    canvas.NewCircle(theme.Color(theme.ColorNamePrimary)),
		initials:  canvas.NewText("", color.White),
		name:      widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		id:        widget.NewLabel(""),
		pill:      canvas.NewRectangle(color.Transparent),
		stateText: canvas.NewText("", color.White),
	}
	r.initials.TextStyle = fyne.TextStyle{Bold: true}
	r.pill.CornerRadius = 10
	r.stateText.TextSize = theme.CaptionTextSize()
	r.stateText.TextStyle = fyne.TextStyle{Bold: true}
	r.ExtendBaseWidget(r)
	return r
}

// CreateRenderer lays out the row: avatar on the left, pill on the right and
// the name over the ID in the middle. Fyne calls it once per row.
func (r *patientRow) CreateRenderer() fyne.WidgetRenderer {
	avatar := container.NewGridWrap(fyne.NewSize(40, 40),
		container.NewStack(r.avatar, container.NewCenter(r.initials)))
	pill := container.NewStack(r.pill, container.NewPadded(r.stateText))
	text := container.NewVBox(r.name, r.id)
	return widget.NewSimpleRenderer(
		container.NewBorder(nil, nil, container.NewCenter(avatar), container.NewCenter(pill), text))
}

// set shows patient p in the row. The final Refresh redraws the row and
// recalculates its layout, because the pill grows or shrinks with its text.
func (r *patientRow) set(p *hospital.Patient) {
	r.initials.Text = initials(p.Name())
	r.name.SetText(p.Name())
	r.id.SetText(p.ID())
	r.pill.FillColor = stateColor(p.State())
	r.stateText.Text = p.State().String()
	r.Refresh()
}

// stateColor returns the color that identifies patient state s in the GUI:
// green when awake, amber when asleep in the hallway and indigo when asleep
// in bed.
func stateColor(s hospital.PatientState) color.Color {
	switch s {
	case hospital.Awake:
		return color.NRGBA{R: 0x16, G: 0xA3, B: 0x4A, A: 0xFF} // green
	case hospital.AsleepInHallway:
		return color.NRGBA{R: 0xD9, G: 0x77, B: 0x06, A: 0xFF} // amber
	case hospital.AsleepInBed:
		return color.NRGBA{R: 0x4F, G: 0x46, B: 0xE5, A: 0xFF} // indigo
	default:
		return color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF} // gray, unknown state
	}
}

// initials returns the first letter of the first two words of name in upper
// case, for example "CP" for "Carlos Polo".
func initials(name string) string {
	words := strings.Fields(name)
	result := ""
	for i := 0; i < len(words) && i < 2; i++ {
		first := []rune(words[i])[0] // Fields never returns empty words
		result += string(unicode.ToUpper(first))
	}
	return result
}

// locationText returns location, or "-" if the patient never fell asleep.
func locationText(location string) string {
	if location == "" {
		return "-"
	}
	return location
}

// roomText describes room r, or says the patient has none.
func roomText(r *hospital.Room) string {
	if r == nil {
		return "Sin cuarto"
	}
	return "Cuarto " + strconv.Itoa(r.Number())
}

// doctorText describes doctor d, or says no doctor is in charge yet.
func doctorText(d *hospital.Doctor) string {
	if d == nil {
		return "Sin asignar"
	}
	return d.Name() + " (" + d.Specialty() + ")"
}

// setEnabled enables button b when enabled is true and disables it otherwise.
func setEnabled(b *widget.Button, enabled bool) {
	if enabled {
		b.Enable()
	} else {
		b.Disable()
	}
}

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// patientsView is the Patients page: the searchable list of patients on the
// left and the detail of the selected patient, with its actions, on the
// right.
type patientsView struct {
	a        *App
	query    string              // text typed in the search box
	patients []*hospital.Patient // patients that match query, shown by list
	selected *hospital.Patient   // nil until the user picks a patient

	list  *widget.List
	count *canvas.Text

	empty  fyne.CanvasObject // shown while no patient is selected
	detail fyne.CanvasObject // shown when a patient is selected

	avatar    *avatar
	name      *widget.Label
	subtitle  *canvas.Text
	statePill *pill
	levelPill *pill
	age       *widget.Label
	location  *widget.Label
	room      *widget.Label
	doctor    *widget.Label

	attackLocation *widget.Entry
	attackBtn      *widget.Button
	wakeBtn        *widget.Button
	assignBtn      *widget.Button
}

// newPatientsPage builds the Patients page of app a and registers its
// refresher.
func newPatientsPage(a *App) page {
	v := &patientsView{a: a}
	return page{
		section:  "Hospital",
		title:    "Pacientes",
		subtitle: "Sin afán, pero sin dormirse: busca, mira y atiende",
		icon:     icon("users"),
		actions: []fyne.CanvasObject{
			newActionButton("Admitir paciente", icon("user-plus"), widget.HighImportance, a.showAdmitPatientForm),
		},
		content: v.build(),
	}
}

// build creates the widgets of the page and returns its content.
func (v *patientsView) build() fyne.CanvasObject {
	search := widget.NewEntry()
	search.SetPlaceHolder("Buscar por nombre o ID")
	search.ActionItem = widget.NewIcon(icon("magnifying-glass"))
	search.OnChanged = func(s string) {
		v.query = s
		// The rows move when the filter changes, so the old highlight would
		// point to another patient. The detail keeps showing v.selected.
		v.list.UnselectAll()
		v.refresh()
	}

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
	v.count = newCaption("")

	left := newPanel(container.NewBorder(
		container.NewVBox(search, padded(4, v.count)), nil, nil, nil, v.list))

	split := container.NewHSplit(left, v.buildDetail())
	split.Offset = 0.42
	v.a.addRefresher(v.refresh)
	return split
}

// buildDetail creates the right panel: an empty state and the detail of the
// selected patient, one on top of the other, so only one is visible.
func (v *patientsView) buildDetail() fyne.CanvasObject {
	v.avatar = newAvatar(72)
	v.name = newHeading("", theme.SizeNameHeadingText)
	v.subtitle = newCaption("")
	v.statePill = newPill()
	v.levelPill = newPill()
	pills := container.NewHBox(v.statePill.box, v.levelPill.box)
	header := container.NewBorder(nil, nil, v.avatar.box, nil,
		padded(6, container.NewVBox(v.name, padded(2, v.subtitle), pills)))

	v.age = widget.NewLabel("")
	v.location = widget.NewLabel("")
	v.room = widget.NewLabel("")
	v.doctor = widget.NewLabel("")
	info := newGroupedList(
		newInfoRow("Edad", v.age),
		newInfoRow("Último ataque en", v.location),
		newInfoRow("Cuarto", v.room),
		newInfoRow("Doctor asignado", v.doctor),
	)

	v.attackLocation = widget.NewEntry()
	v.attackLocation.SetPlaceHolder("¿Dónde se durmió? p. ej. cafeteria")
	v.attackBtn = newActionButton("Ataque de sueño", icon("moon-stars"), widget.HighImportance, v.sleepAttack)
	v.wakeBtn = newActionButton("Despertar", icon("sun-horizon"), widget.MediumImportance, v.wakeUp)
	v.assignBtn = newActionButton("Asignar cuarto", icon("bed"), widget.HighImportance, v.assignRoom)
	actions := container.NewVBox(
		newHeading("Acciones", theme.SizeNameSubHeadingText),
		container.NewBorder(nil, nil, nil, v.attackBtn, v.attackLocation),
		container.NewGridWithColumns(2, v.wakeBtn, v.assignBtn),
	)

	v.detail = container.NewVBox(header, padded(6, info), actions)

	hint := widget.NewLabelWithStyle("Escoge un paciente de la lista, mi llave", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	v.empty = container.NewCenter(container.NewVBox(
		container.NewCenter(newSombrero(120, 78)), hint,
		container.NewCenter(newCaption("Aquí vas a ver su estado y lo que puedes hacer"))))

	return newPanel(container.NewStack(v.empty, container.NewVScroll(v.detail)))
}

// refresh reloads the list and the detail from the hospital. App.refreshAll
// calls it after every action.
func (v *patientsView) refresh() {
	all := v.a.h.Patients()
	v.patients = filterPatients(all, v.query)
	v.count.Text = fmt.Sprintf("Mostrando %d de %d pacientes", len(v.patients), len(all))
	v.count.Refresh()
	v.list.Refresh()
	v.refreshDetail()
}

// filterPatients returns the patients whose name or ID contain query,
// ignoring upper and lower case. An empty query keeps every patient.
func filterPatients(patients []*hospital.Patient, query string) []*hospital.Patient {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return patients
	}
	var found []*hospital.Patient
	for _, p := range patients {
		if strings.Contains(strings.ToLower(p.Name()+" "+p.ID()), query) {
			found = append(found, p)
		}
	}
	return found
}

// refreshDetail shows the selected patient and enables only the actions
// that make sense for their current state.
func (v *patientsView) refreshDetail() {
	p := v.selected
	if p == nil {
		v.detail.Hide()
		v.empty.Show()
		return
	}
	v.empty.Hide()
	v.detail.Show()

	v.avatar.set(p.Name())
	v.name.SetText(p.Name())
	v.subtitle.Text = "Paciente " + p.ID()
	v.subtitle.Refresh()
	v.statePill.set(stateLabel(p.State()), stateColor(p.State()))
	v.levelPill.set("Narcolepsia "+strings.ToLower(levelLabel(p.Level())), levelColor(p.Level()))
	v.age.SetText(strconv.Itoa(p.Age()) + " años")
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
		v.a.notify(fmt.Sprintf("%s ya tiene cama: cuarto %d", p.Name(), p.Room().Number()))
	case errors.Is(err, hospital.ErrNoRoomAvailable):
		// The episode is registered anyway: the patient sleeps in the hallway.
		v.a.lastAssign = &assignResult{patient: p, err: err}
		v.a.showWarning("¡Erda! Episodio registrado, pero no hay cuarto: el paciente queda en el pasillo")
	default:
		// Nothing was registered, so it is not an assignment result.
		v.a.showError(err)
	}
	v.attackLocation.SetText("")
	v.a.refreshAll()
}

// wakeUp wakes the selected patient up, releasing their room if they had one.
func (v *patientsView) wakeUp() {
	if err := v.a.h.WakeUpPatient(v.selected); err != nil {
		v.a.showError(err)
	} else {
		v.a.notify(v.selected.Name() + " se despertó. ¡Buenos días, mi llave!")
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
	case err == nil:
		v.a.notify(fmt.Sprintf("%s pasó al cuarto %d", p.Name(), r.Number()))
	case errors.Is(err, hospital.ErrNoRoomAvailable):
		v.a.showWarning("Ombe, no hay cuartos libres: el paciente sigue en el pasillo")
	default:
		v.a.showError(err)
	}
	v.a.refreshAll()
}

// patientRow is one row of the patient list: an avatar with initials, the
// name over the ID, age and level, and a pill colored by state. The list
// creates only a few rows and reuses them, so set changes what a row shows.
type patientRow struct {
	widget.BaseWidget
	avatar *avatar
	name   *canvas.Text
	info   *canvas.Text
	state  *pill
}

// newPatientRow returns an empty row, ready to be filled with set.
func newPatientRow() *patientRow {
	r := &patientRow{
		avatar: newAvatar(40),
		name:   canvas.NewText("", ink),
		info:   newCaption(""),
		state:  newPill(),
	}
	r.name.TextStyle = fyne.TextStyle{Bold: true}
	r.ExtendBaseWidget(r)
	return r
}

// CreateRenderer lays out the row: avatar on the left, pill on the right and
// the name over the details in the middle. Fyne calls it once per row.
func (r *patientRow) CreateRenderer() fyne.WidgetRenderer {
	text := container.NewCenter(container.NewVBox(r.name, r.info))
	return widget.NewSimpleRenderer(padded(4, container.NewBorder(nil, nil,
		r.avatar.box, r.state.box, container.NewHBox(padded(2, text)))))
}

// set shows patient p in the row. The final Refresh redraws the row and
// recalculates its layout, because the texts change their width.
func (r *patientRow) set(p *hospital.Patient) {
	r.avatar.set(p.Name())
	r.name.Text = p.Name()
	r.info.Text = fmt.Sprintf("%s · %d años · %s", p.ID(), p.Age(), levelLabel(p.Level()))
	r.state.set(stateLabel(p.State()), stateColor(p.State()))
	r.Refresh()
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

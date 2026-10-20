package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// queriesView is the Queries tab: the four queries of the assignment, one
// per inner tab. Inner tabs are used instead of an accordion because a list
// or a table inside an accordion only gets the height of one row.
type queriesView struct {
	a *App

	hallway     []string // lines shown by hallwayList
	hallwayList *widget.List

	doctors      []*hospital.Doctor
	doctorSelect *widget.Select
	episodes     []string // lines shown by episodesList
	episodesList *widget.List

	assignment *widget.Label

	report      map[*hospital.Patient]int
	severe      []*hospital.Patient // keys of report, sorted by ID
	reportTable *widget.Table
}

// reportHeaders are the column titles of the severity report table.
var reportHeaders = []string{"ID", "Paciente", "Episodios hoy"}

// newQueriesTab builds the Queries tab of app a and registers its refresher.
func newQueriesTab(a *App) fyne.CanvasObject {
	v := &queriesView{a: a, doctors: a.h.Doctors()}

	v.hallwayList = newTextList(func() []string { return v.hallway })

	names := make([]string, 0, len(v.doctors))
	for _, d := range v.doctors {
		names = append(names, d.Name()+" ("+d.Specialty()+")")
	}
	v.episodesList = newTextList(func() []string { return v.episodes })
	v.doctorSelect = widget.NewSelect(names, func(string) { v.refreshEpisodes() })
	v.doctorSelect.PlaceHolder = "Elige un doctor"

	v.assignment = widget.NewLabel("")
	v.assignment.Wrapping = fyne.TextWrapWord

	v.reportTable = widget.NewTable(
		func() (int, int) {
			return len(v.severe), len(reportHeaders)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			p := v.severe[id.Row]
			cell := obj.(*widget.Label)
			switch id.Col {
			case 0:
				cell.SetText(p.ID())
			case 1:
				cell.SetText(p.Name())
			default:
				cell.SetText(strconv.Itoa(v.report[p]))
			}
		},
	)
	v.reportTable.ShowHeaderRow = true
	v.reportTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		obj.(*widget.Label).SetText(reportHeaders[id.Col])
	}
	v.reportTable.SetColumnWidth(0, 100)
	v.reportTable.SetColumnWidth(1, 240)
	v.reportTable.SetColumnWidth(2, 140)

	if len(v.doctors) > 0 {
		v.doctorSelect.SetSelectedIndex(0)
	}
	a.addRefresher(v.refresh)

	return container.NewAppTabs(
		container.NewTabItem("Pacientes en el pasillo", v.hallwayList),
		container.NewTabItem("Episodios por doctor",
			container.NewBorder(v.doctorSelect, nil, nil, nil, v.episodesList)),
		container.NewTabItem("Asignación de cuarto", container.NewVBox(v.assignment)),
		container.NewTabItem("Reporte de severidad", v.reportTable),
	)
}

// refresh runs the four queries again on the hospital. App.refreshAll calls
// it after every action.
func (v *queriesView) refresh() {
	var hallway []string
	for _, p := range v.a.h.PatientsInHallway() {
		hallway = append(hallway, fmt.Sprintf("%s  %s (%s)", p.ID(), p.Name(), p.Location()))
	}
	v.hallway = orNone(hallway)
	v.hallwayList.Refresh()

	v.refreshEpisodes()

	v.assignment.SetText(assignmentText(v.a.lastAssign))

	v.report = v.a.h.SevereReport()
	v.severe = sortedByID(v.report)
	v.reportTable.Refresh()
}

// refreshEpisodes shows the episodes of the doctor chosen in the select.
func (v *queriesView) refreshEpisodes() {
	i := v.doctorSelect.SelectedIndex()
	if i < 0 {
		v.episodes = []string{"(elige un doctor)"}
	} else {
		v.episodes = orNone(summaries(v.doctors[i].MyEpisodes()))
	}
	v.episodesList.Refresh()
}

// assignmentText describes the last room assignment r, which is nil if no
// assignment has been made yet.
func assignmentText(r *assignResult) string {
	switch {
	case r == nil:
		return "Todavía no se ha asignado ningún cuarto."
	case r.err == nil:
		return fmt.Sprintf("%s %s -> cuarto %d", r.patient.ID(), r.patient.Name(), r.room.Number())
	case errors.Is(r.err, hospital.ErrNoRoomAvailable):
		return fmt.Sprintf("%s %s: no hay cuarto disponible, sigue en el pasillo.\n(%v)",
			r.patient.ID(), r.patient.Name(), r.err)
	default:
		return fmt.Sprintf("%s %s: error: %v", r.patient.ID(), r.patient.Name(), r.err)
	}
}

// sortedByID returns the patients of report sorted by ID. A map has no
// order, and Go even changes it on purpose between iterations, so without
// sorting the rows of the table would jump around on every refresh.
func sortedByID(report map[*hospital.Patient]int) []*hospital.Patient {
	patients := make([]*hospital.Patient, 0, len(report))
	for p := range report {
		patients = append(patients, p)
	}
	slices.SortFunc(patients, func(a, b *hospital.Patient) int {
		return cmp.Compare(a.ID(), b.ID())
	})
	return patients
}

// newTextList returns a widget.List with one label per line returned by
// lines. Change the data behind lines and call Refresh on the list to show it.
func newTextList(lines func() []string) *widget.List {
	return widget.NewList(
		func() int {
			return len(lines())
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(lines()[i])
		},
	)
}

// summaries returns the Summary of every record, in the same order.
func summaries(records []hospital.EpisodeRecord) []string {
	lines := make([]string, 0, len(records))
	for _, rec := range records {
		lines = append(lines, rec.Summary())
	}
	return lines
}

// orNone returns lines, or a single "(ninguno)" line if lines is empty, so
// an empty query still shows something.
func orNone(lines []string) []string {
	if len(lines) == 0 {
		return []string{"(ninguno)"}
	}
	return lines
}

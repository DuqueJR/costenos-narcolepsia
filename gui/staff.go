package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// newStaffPage builds the Staff page of app a: one card per doctor and one
// per orderly, with how many episodes each one attended. Staff can be hired
// from the GUI, so the refresher builds the cards again every time.
func newStaffPage(a *App) page {
	doctors := container.NewGridWithColumns(3)
	orderlies := container.NewGridWithColumns(3)

	a.addRefresher(func() {
		attended := attendedByStaff(a.h.History())

		doctors.RemoveAll()
		isDoctor := make(map[string]bool)
		for _, d := range a.h.Doctors() {
			isDoctor[d.ID()] = true
			doctors.Add(newStaffCard(d.Name(), d.ID(), d.Specialty(), []stat{
				{"Pacientes a cargo", len(d.Patients())},
				{"Episodios atendidos", attended[d.ID()]},
			}))
		}
		doctors.Refresh()

		// Staff holds every Attender. The only ones the GUI can hire besides
		// doctors are orderlies, so whoever is not a doctor is shown here.
		// It only uses ID and Name, the methods of the Attender interface,
		// without asking for the concrete type.
		orderlies.RemoveAll()
		for _, s := range a.h.Staff() {
			if !isDoctor[s.ID()] {
				orderlies.Add(newStaffCard(s.Name(), s.ID(), "Camillero", []stat{
					{"Episodios atendidos", attended[s.ID()]},
				}))
			}
		}
		orderlies.Refresh()
	})

	content := container.NewVBox(
		newHeading("Doctores", theme.SizeNameSubHeadingText), doctors,
		newHeading("Camilleros", theme.SizeNameSubHeadingText), orderlies,
	)
	return page{
		title:    "Personal",
		subtitle: "Quién atiende las emergencias, por turnos",
		icon:     theme.LoginIcon(),
		actions: []fyne.CanvasObject{
			newActionButton("Contratar doctor", theme.ContentAddIcon(), widget.HighImportance, a.showHireDoctorForm),
			newActionButton("Contratar camillero", theme.ContentAddIcon(), widget.MediumImportance, a.showHireOrderlyForm),
		},
		content: container.NewVScroll(content),
	}
}

// stat is one number shown in a staff card, with its caption.
type stat struct {
	caption string
	value   int
}

// newStaffCard builds the card of a staff member: avatar, name, ID, role
// and a row of numbers.
func newStaffCard(name, id, role string, stats []stat) fyne.CanvasObject {
	face := newAvatar(52)
	face.set(name, id)
	title := widget.NewLabelWithStyle(name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewBorder(nil, nil, face.box, nil,
		container.NewVBox(title, padded(4, newCaption(id+" · "+role))))

	numbers := container.NewGridWithColumns(len(stats))
	for _, s := range stats {
		value := newHeading(strconv.Itoa(s.value), theme.SizeNameSubHeadingText)
		numbers.Add(container.NewVBox(value, padded(4, newCaption(s.caption))))
	}
	return newPanel(container.NewVBox(header, widget.NewSeparator(), numbers))
}

// attendedByStaff counts, from history, how many episodes each staff member
// attended, keyed by staff ID.
func attendedByStaff(history []hospital.EpisodeRecord) map[string]int {
	counts := make(map[string]int)
	for _, rec := range history {
		counts[rec.AttendedBy().ID()]++
	}
	return counts
}

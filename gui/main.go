// Command gui is the Fyne desktop interface of the hospital. It shows the
// same hospital package used by the console program, without changing it.
package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

func main() {
	h, err := newSeedHospital()
	if err != nil {
		// There is no window yet, so the error goes to the console.
		log.Fatalf("building the hospital: %v", err)
	}

	fyneApp := app.NewWithID("co.edu.eia.costenos")
	win := fyneApp.NewWindow("Hospital de los Costeños con Narcolepsia")
	win.Resize(fyne.NewSize(1100, 720))

	a := newApp(h, win)
	header := newHeader(a)
	tabs := newTabs(a)
	win.SetContent(container.NewBorder(header, nil, nil, nil, tabs))

	a.refreshAll()
	win.ShowAndRun()
}

// newHeader builds the row of counters shown above the tabs: patients,
// patients asleep in the hallway, patients asleep in bed and free rooms.
// It registers a refresher that recalculates them from the hospital getters.
func newHeader(a *App) fyne.CanvasObject {
	patients := widget.NewCard("", "Pacientes", nil)
	inHallway := widget.NewCard("", "Dormidos en pasillo", nil)
	inBed := widget.NewCard("", "Dormidos en cama", nil)
	freeRooms := widget.NewCard("", "Cuartos libres", nil)

	a.addRefresher(func() {
		asleepInBed := 0
		for _, p := range a.h.Patients() {
			if p.State() == hospital.AsleepInBed {
				asleepInBed++
			}
		}
		free := 0
		for _, r := range a.h.Rooms() {
			if r.IsAvailable() {
				free++
			}
		}
		patients.SetTitle(strconv.Itoa(len(a.h.Patients())))
		inHallway.SetTitle(strconv.Itoa(len(a.h.PatientsInHallway())))
		inBed.SetTitle(strconv.Itoa(asleepInBed))
		freeRooms.SetTitle(strconv.Itoa(free))
	})

	return container.NewGridWithColumns(4, patients, inHallway, inBed, freeRooms)
}

// newTabs builds the tabs of the main window of app a, each one with its
// theme icon. The tabs that are not built yet show a placeholder.
func newTabs(a *App) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItemWithIcon("Pacientes", theme.AccountIcon(), newPatientsTab(a)),
		container.NewTabItemWithIcon("Cuartos", theme.HomeIcon(), newRoomsTab(a)),
		container.NewTabItemWithIcon("Consultas", theme.SearchIcon(), newQueriesTab(a)),
		container.NewTabItemWithIcon("Historial", theme.HistoryIcon(), comingSoon()),
		container.NewTabItemWithIcon("Música", theme.MediaMusicIcon(), comingSoon()),
	)
}

// comingSoon returns a centered "Próximamente" label for a tab that is not
// built yet.
func comingSoon() fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("Próximamente"))
}

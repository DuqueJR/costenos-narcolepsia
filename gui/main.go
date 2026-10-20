// Command gui is the Fyne desktop interface of the hospital. It shows the
// same hospital package used by the console program, without changing it.
package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	h, err := newSeedHospital()
	if err != nil {
		// There is no window yet, so the error goes to the console.
		log.Fatalf("building the hospital: %v", err)
	}

	fyneApp := app.NewWithID("co.edu.eia.costenos")
	fyneApp.Settings().SetTheme(newCaribeTheme())
	win := fyneApp.NewWindow("Hospital de los Costeños con Narcolepsia")
	win.Resize(fyne.NewSize(1280, 800))

	a := newApp(h, win)
	win.SetContent(newShell(a, newPages(a)))

	a.refreshAll()
	win.ShowAndRun()
}

// newPages builds every page of the sidebar for app a, in menu order.
func newPages(a *App) []page {
	return []page{
		newDashboardPage(a),
		newPatientsPage(a),
		newRoomsPage(a),
		newStaffPage(a),
		newQueriesPage(a),
		newHistoryPage(a),
		newMusicPage(),
	}
}

// Command gui is the Fyne desktop interface of the hospital. It shows the
// same hospital package used by the console program, without changing it.
package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	h, err := newSeedHospital()
	if err != nil {
		// There is no window yet, so the error goes to the console.
		log.Fatalf("building the hospital: %v", err)
	}

	fyneApp := app.NewWithID("co.edu.eia.costenos")
	fyneApp.Settings().SetTheme(caribeTheme{})
	win := fyneApp.NewWindow("Hospital de los Costeños con Narcolepsia")
	win.Resize(fyne.NewSize(1240, 780))

	a := newApp(h, win)
	win.SetContent(newShell(a, newPages(a)))

	a.refreshAll()
	win.ShowAndRun()
}

// newPages builds every page of the sidebar for app a, in menu order.
func newPages(a *App) []page {
	return []page{
		newPatientsPage(a),
		newRoomsPage(a),
		{title: "Consultas", subtitle: "Las cuatro consultas del enunciado", icon: theme.SearchIcon(), content: newQueriesTab(a)},
		{title: "Historial", subtitle: "Episodios del más reciente al más antiguo", icon: theme.HistoryIcon(), content: newHistoryTab(a)},
		newMusicPage(),
	}
}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"costenos-narcolepsia/hospital"
)

// assignResult is the outcome of the last room assignment made from the GUI,
// shown later by the "Room assignment" query.
type assignResult struct {
	patient *hospital.Patient
	room    *hospital.Room // nil when the assignment failed
	err     error
}

// App holds the state shared by every view of the GUI: the hospital, the main
// window, the last room assignment and the functions that refresh each view.
type App struct {
	h          *hospital.Hospital
	win        fyne.Window
	lastAssign *assignResult // nil until the first assignment
	refreshers []func()      // one per view, plus the header
}

// newApp returns an App that shows hospital h in window win.
func newApp(h *hospital.Hospital, win fyne.Window) *App {
	return &App{h: h, win: win}
}

// addRefresher registers f so that refreshAll calls it. Each view registers
// its own function when it is built, so App does not need to know them.
func (a *App) addRefresher(f func()) {
	a.refreshers = append(a.refreshers, f)
}

// refreshAll redraws every view with the current state of the hospital. It is
// called after every action, from the UI goroutine.
func (a *App) refreshAll() {
	for _, refresh := range a.refreshers {
		refresh()
	}
}

// showError shows err in an error dialog over the main window.
func (a *App) showError(err error) {
	dialog.ShowError(err, a.win)
}

// showWarning shows msg in an information dialog over the main window. It is
// meant for expected situations, such as a patient left in the hallway.
func (a *App) showWarning(msg string) {
	dialog.ShowInformation("Aviso", msg, a.win)
}

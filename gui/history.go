package main

import (
	"slices"

	"fyne.io/fyne/v2"
)

// newHistoryTab builds the History tab of app a: the Summary of every
// episode of the hospital, from the newest to the oldest.
func newHistoryTab(a *App) fyne.CanvasObject {
	var lines []string
	list := newTextList(func() []string { return lines })

	a.addRefresher(func() {
		// History returns a copy, so reversing it does not change the hospital.
		history := a.h.History()
		slices.Reverse(history)
		lines = orNone(summaries(history))
		list.Refresh()
	})
	return list
}

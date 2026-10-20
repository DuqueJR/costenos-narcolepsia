package main

import (
	"slices"

	"fyne.io/fyne/v2/theme"
)

// newHistoryPage builds the History page of app a: the Summary of every
// episode of the hospital, from the newest to the oldest.
func newHistoryPage(a *App) page {
	var lines []string
	list := newTextList(func() []string { return lines })

	a.addRefresher(func() {
		// History returns a copy, so reversing it does not change the hospital.
		history := a.h.History()
		slices.Reverse(history)
		lines = orNone(summaries(history))
		list.Refresh()
	})
	return page{
		title:    "Historial",
		subtitle: "Todos los episodios, del más reciente al más antiguo",
		icon:     theme.HistoryIcon(),
		content:  newPanel(list),
	}
}

package main

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// newRoomsTab builds the Rooms tab of app a: a grid of 3 columns with one
// card per room. The GUI never adds rooms, so the cards are built once and
// each one registers a refresher that updates its state and occupants.
func newRoomsTab(a *App) fyne.CanvasObject {
	grid := container.NewGridWithColumns(3)
	for _, r := range a.h.Rooms() {
		grid.Add(newRoomCard(a, r))
	}
	return container.NewVScroll(grid)
}

// newRoomCard builds the card of room r and registers its refresher in a.
func newRoomCard(a *App, r *hospital.Room) fyne.CanvasObject {
	capacity := widget.NewLabel("Capacidad: " + strconv.Itoa(r.Capacity()) + " cama(s)")
	state := canvas.NewText("", nil)
	state.TextStyle = fyne.TextStyle{Bold: true}
	occupants := widget.NewLabel("")

	a.addRefresher(func() {
		state.Text = r.State().String()
		state.Color = roomStateColor(r.State())
		state.Refresh()
		occupants.SetText(occupantsText(r.Occupants()))
	})

	title := "Cuarto " + strconv.Itoa(r.Number())
	return widget.NewCard(title, "", container.NewVBox(capacity, state, occupants))
}

// occupantsText returns the names of the patients in a room, or "Libre" if
// it is empty.
func occupantsText(occupants []*hospital.Patient) string {
	if len(occupants) == 0 {
		return "Libre"
	}
	names := make([]string, 0, len(occupants))
	for _, p := range occupants {
		names = append(names, p.Name())
	}
	return "Ocupante: " + strings.Join(names, ", ")
}

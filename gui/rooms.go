package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// newRoomsPage builds the Rooms page of app a: a grid of 3 columns with one
// card per room. Rooms can be added from the GUI, so the refresher builds
// the cards again every time; with a few rooms that is simpler than keeping
// each card and updating it.
func newRoomsPage(a *App) page {
	grid := container.NewGridWithColumns(3)
	a.addRefresher(func() {
		grid.RemoveAll()
		for _, r := range a.h.Rooms() {
			grid.Add(newRoomCard(r))
		}
		grid.Refresh()
	})
	return page{
		section:  "Hospital",
		title:    "Cuartos",
		subtitle: "Camas listas pa' cuando a alguien le pique el sueño",
		icon:     theme.HomeIcon(),
		color:    green,
		actions: []fyne.CanvasObject{
			newActionButton("Agregar cuarto", theme.ContentAddIcon(), widget.HighImportance, a.showAddRoomForm),
		},
		content: container.NewVScroll(container.NewVBox(grid)),
	}
}

// newRoomCard builds the card of room r: its number and state, a bar with
// the beds in use and the patients sleeping in it.
func newRoomCard(r *hospital.Room) fyne.CanvasObject {
	occupants := r.Occupants()

	state := newPill()
	state.set(roomStateLabel(r.State()), roomStateColor(r.State()))
	title := newHeading("Cuarto "+strconv.Itoa(r.Number()), theme.SizeNameSubHeadingText)

	beds := widget.NewProgressBar()
	beds.Max = float64(r.Capacity())
	beds.TextFormatter = func() string {
		return fmt.Sprintf("%d de %d camas ocupadas", len(occupants), r.Capacity())
	}
	beds.SetValue(float64(len(occupants)))

	people := container.NewVBox()
	if len(occupants) == 0 {
		people.Add(padded(6, newCaption("Libre, esperando al próximo dormilón")))
	}
	for _, p := range occupants {
		face := newAvatar(30)
		face.set(p.Name(), p.ID())
		name := widget.NewLabelWithStyle(p.Name(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		people.Add(container.NewBorder(nil, nil, face.box, nil, container.NewHBox(name, container.NewCenter(newCaption(p.ID())))))
	}

	return newPanel(container.NewVBox(
		container.NewBorder(nil, nil, newIconBadge(theme.HomeIcon(), roomStateColor(r.State()), 30), state.box, title),
		padded(2, beds),
		people,
	))
}

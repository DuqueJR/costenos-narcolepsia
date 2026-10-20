package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// page is one screen of the sidebar: its title and subtitle, its icon in the
// menu, the buttons shown next to the title and its content.
type page struct {
	title    string
	subtitle string
	icon     fyne.Resource
	actions  []fyne.CanvasObject
	content  fyne.CanvasObject
}

// newShell builds the frame of the window for app a: the sidebar with one
// button per page on the left and, on the right, the title of the current
// page, the stats of the hospital and the content of the page.
func newShell(a *App, pages []page) fyne.CanvasObject {
	title := newHeading("", theme.SizeNameHeadingText)
	subtitle := newCaption("")
	actions := container.NewHBox()
	body := container.NewStack()
	buttons := make([]*widget.Button, len(pages))

	// show puts page i on screen and highlights its button in the menu.
	show := func(i int) {
		p := pages[i]
		title.SetText(p.title)
		subtitle.Text = p.subtitle
		subtitle.Refresh()
		actions.Objects = p.actions
		actions.Refresh()
		body.Objects = []fyne.CanvasObject{p.content}
		body.Refresh()
		for j, b := range buttons {
			b.Importance = widget.LowImportance
			if j == i {
				b.Importance = widget.HighImportance
			}
			b.Refresh()
		}
	}

	menu := container.NewVBox()
	for i, p := range pages {
		b := widget.NewButtonWithIcon(p.title, p.icon, func() { show(i) })
		b.Alignment = widget.ButtonAlignLeading
		buttons[i] = b
		menu.Add(b)
	}
	show(0)

	heading := container.NewBorder(nil, nil, nil, container.NewCenter(actions),
		container.NewVBox(title, padded(4, subtitle)))
	top := container.NewVBox(heading, newStatsStrip(a))
	right := padded(18, container.NewBorder(top, nil, nil, nil, body))
	return container.NewBorder(nil, nil, newSidebar(menu), nil, right)
}

// newSidebar returns the white column on the left: the name of the hospital
// on top and the menu below it.
func newSidebar(menu fyne.CanvasObject) fyne.CanvasObject {
	logoBg := canvas.NewCircle(turquoise)
	logoIcon := widget.NewIcon(theme.NewColoredResource(theme.ContentAddIcon(), theme.ColorNameForegroundOnPrimary))
	logo := container.NewGridWrap(fyne.NewSize(42, 42), container.NewStack(logoBg, padded(8, logoIcon)))

	name := canvas.NewText("Costeños", ink)
	name.TextSize = 18
	name.TextStyle = fyne.TextStyle{Bold: true}
	brand := container.NewBorder(nil, nil, container.NewCenter(logo), nil,
		padded(4, container.NewVBox(name, newCaption("con Narcolepsia"))))

	section := newCaption("MENÚ")
	section.TextStyle = fyne.TextStyle{Bold: true}
	content := container.NewBorder(container.NewVBox(brand, widget.NewSeparator(), padded(6, section)),
		nil, nil, nil, menu)

	bg := canvas.NewRectangle(white)
	width := canvas.NewRectangle(white)
	width.SetMinSize(fyne.NewSize(210, 0)) // fixed width of the sidebar
	line := canvas.NewRectangle(border)
	line.SetMinSize(fyne.NewSize(1, 0))
	return container.NewBorder(nil, nil, nil, line, container.NewStack(bg, width, padded(14, content)))
}

// newStatsStrip builds the row of counters shown on every page: patients,
// patients asleep in the hallway, patients asleep in bed and free rooms. It
// registers a refresher that recalculates them from the hospital getters.
func newStatsStrip(a *App) fyne.CanvasObject {
	patients := newStatCard("Pacientes", theme.AccountIcon(), turquoise)
	inHallway := newStatCard("Dormidos en pasillo", theme.WarningIcon(), amber)
	inBed := newStatCard("Dormidos en cama", theme.VisibilityOffIcon(), indigo)
	freeRooms := newStatCard("Cuartos libres", theme.HomeIcon(), green)

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
		patients.set(len(a.h.Patients()))
		inHallway.set(len(a.h.PatientsInHallway()))
		inBed.set(asleepInBed)
		freeRooms.set(free)
	})

	strip := container.NewGridWithColumns(4, patients.box, inHallway.box, inBed.box, freeRooms.box)
	return container.New(layout.NewCustomPaddedLayout(12, 12, 0, 0), strip)
}

// newMusicPage is the Music page. The player arrives in another branch, so
// for now it only says so.
func newMusicPage() page {
	icon := widget.NewIcon(theme.NewColoredResource(theme.MediaMusicIcon(), theme.ColorNamePrimary))
	big := container.NewGridWrap(fyne.NewSize(72, 72), icon)
	text := newHeading("Próximamente", theme.SizeNameSubHeadingText)
	text.Alignment = fyne.TextAlignCenter
	return page{
		title:    "Música",
		subtitle: "Un reproductor para la sala de espera",
		icon:     theme.MediaMusicIcon(),
		content:  newPanel(container.NewCenter(container.NewVBox(container.NewCenter(big), text))),
	}
}

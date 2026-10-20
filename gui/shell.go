package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"costenos-narcolepsia/hospital"
)

// page is one screen of the sidebar: the section of the menu it belongs to,
// its title and subtitle, its icon and color in the menu, the buttons shown
// next to the title and its content.
type page struct {
	section  string
	title    string
	subtitle string
	icon     fyne.Resource
	color    color.Color
	actions  []fyne.CanvasObject
	content  fyne.CanvasObject
}

// newShell builds the frame of the window for app a: the sidebar with one
// item per page on the left and, on the right, the title of the current
// page, the stats of the hospital and the content of the page.
func newShell(a *App, pages []page) fyne.CanvasObject {
	title := newHeading("", theme.SizeNameHeadingText)
	subtitle := newCaption("")
	subtitle.TextSize = 13
	actions := container.NewHBox()
	body := container.NewStack()
	items := make([]*sidebarItem, len(pages))

	// show puts page i on screen and highlights its item in the menu.
	show := func(i int) {
		p := pages[i]
		title.SetText(p.title)
		subtitle.Text = p.subtitle
		subtitle.Refresh()
		actions.Objects = p.actions
		actions.Refresh()
		body.Objects = []fyne.CanvasObject{p.content}
		body.Refresh()
		for j, item := range items {
			item.setSelected(j == i)
		}
	}

	menu := container.NewVBox()
	for i, p := range pages {
		if i == 0 || p.section != pages[i-1].section {
			menu.Add(newSectionTitle(p.section))
		}
		items[i] = newSidebarItem(p.title, p.icon, p.color, func() { show(i) })
		menu.Add(items[i])
	}
	show(0)

	heading := container.NewBorder(nil, nil, nil, container.NewCenter(actions),
		container.NewVBox(title, container.New(layout.NewCustomPaddedLayout(0, 0, 8, 0), subtitle)))
	top := container.NewVBox(heading, newStatsStrip(a))
	right := container.New(layout.NewCustomPaddedLayout(20, 20, 28, 28),
		container.NewBorder(top, nil, nil, nil, body))
	return container.NewBorder(nil, nil, newSidebar(menu), nil, right)
}

// newSidebar returns the gray column on the left: the app icon and name on
// top, the menu in the middle and a costeño greeting at the bottom.
func newSidebar(menu fyne.CanvasObject) fyne.CanvasObject {
	name := canvas.NewText("Costeños", ink)
	name.TextSize = 16
	name.TextStyle = fyne.TextStyle{Bold: true}
	brand := container.NewBorder(nil, nil, newAppIcon(40), nil,
		container.New(layout.NewCustomPaddedLayout(2, 0, 10, 0), container.NewVBox(name, newCaption("Hospital · Narcolepsia"))))

	hello := canvas.NewText("¡Ajá, mi llave!", ink)
	hello.TextStyle = fyne.TextStyle{Bold: true}
	footer := newPanel(container.NewVBox(
		newVueltiaoBand(10, 12), padded(2, hello),
		newCaption("Aquí nadie se queda dormío"), newCaption("sin que lo atiendan."),
	))

	content := container.NewBorder(padded(4, brand), footer, nil, nil, container.NewVScroll(menu))

	bg := canvas.NewRectangle(sidebarBg)
	width := canvas.NewRectangle(color.Transparent)
	width.SetMinSize(fyne.NewSize(236, 0)) // fixed width of the sidebar
	line := canvas.NewRectangle(hairline)
	line.SetMinSize(fyne.NewSize(1, 0))
	return container.NewBorder(nil, nil, nil, line, container.NewStack(bg, width, padded(14, content)))
}

// sidebarItem is one entry of the menu: an icon in a colored square and a
// title. It is a small widget of its own because it reacts to the mouse:
// tapping it runs onTap and hovering it shades its background.
type sidebarItem struct {
	widget.BaseWidget
	bg       *canvas.Rectangle
	label    *canvas.Text
	badge    fyne.CanvasObject
	onTap    func()
	selected bool
}

// Compile-time checks: the build fails if sidebarItem stops reacting to
// taps, to the mouse or stops showing the hand cursor.
var (
	_ fyne.Tappable      = (*sidebarItem)(nil)
	_ desktop.Hoverable  = (*sidebarItem)(nil)
	_ desktop.Cursorable = (*sidebarItem)(nil)
)

// newSidebarItem returns a menu item with title and icon over color c that
// runs onTap when tapped.
func newSidebarItem(title string, icon fyne.Resource, c color.Color, onTap func()) *sidebarItem {
	s := &sidebarItem{
		bg:    canvas.NewRectangle(color.Transparent),
		label: canvas.NewText(title, ink),
		badge: newIconBadge(icon, c, 24),
		onTap: onTap,
	}
	s.bg.CornerRadius = 8
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer lays out the item: the badge on the left and the title.
func (s *sidebarItem) CreateRenderer() fyne.WidgetRenderer {
	row := container.NewBorder(nil, nil, container.NewCenter(s.badge), nil,
		container.New(layout.NewCustomPaddedLayout(0, 0, 10, 0), container.NewVBox(layout.NewSpacer(), s.label, layout.NewSpacer())))
	return widget.NewSimpleRenderer(container.NewStack(s.bg,
		container.New(layout.NewCustomPaddedLayout(6, 6, 8, 8), row)))
}

// setSelected highlights the item with the main color when on is true.
func (s *sidebarItem) setSelected(on bool) {
	s.selected = on
	s.bg.FillColor = color.Transparent
	s.label.Color = ink
	s.label.TextStyle = fyne.TextStyle{}
	if on {
		s.bg.FillColor = caribe
		s.label.Color = white
		s.label.TextStyle = fyne.TextStyle{Bold: true}
	}
	s.Refresh()
}

// Tapped runs the action of the item.
func (s *sidebarItem) Tapped(*fyne.PointEvent) {
	s.onTap()
}

// MouseIn shades the item while the mouse is over it.
func (s *sidebarItem) MouseIn(*desktop.MouseEvent) {
	if !s.selected {
		s.bg.FillColor = color.NRGBA{A: 0x0F}
		s.bg.Refresh()
	}
}

// MouseMoved does nothing; it is part of desktop.Hoverable.
func (s *sidebarItem) MouseMoved(*desktop.MouseEvent) {}

// MouseOut removes the shade when the mouse leaves the item.
func (s *sidebarItem) MouseOut() {
	if !s.selected {
		s.bg.FillColor = color.Transparent
		s.bg.Refresh()
	}
}

// Cursor shows a hand over the item, as over a link.
func (s *sidebarItem) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// newStatsStrip builds the row of counters shown on every page: patients,
// patients asleep in the hallway, patients asleep in bed and free rooms. It
// registers a refresher that recalculates them from the hospital getters.
func newStatsStrip(a *App) fyne.CanvasObject {
	patients := newStatCard("Pacientes", theme.AccountIcon(), caribe)
	inHallway := newStatCard("Dormidos en pasillo", theme.WarningIcon(), orange)
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
	return container.New(layout.NewCustomPaddedLayout(16, 16, 0, 0), strip)
}

// newMusicPage is the Music page. The player arrives in another branch, so
// for now it only says so.
func newMusicPage() page {
	text := newHeading("Próximamente", theme.SizeNameSubHeadingText)
	text.Alignment = fyne.TextAlignCenter
	hint := newCaption("Champeta, vallenato y cumbia pa' la sala de espera")
	return page{
		section:  "Extras",
		title:    "Música",
		subtitle: "Un reproductor para la sala de espera",
		icon:     theme.MediaMusicIcon(),
		color:    pink,
		content: newPanel(container.NewCenter(container.NewVBox(
			container.NewCenter(newIconBadge(theme.MediaMusicIcon(), pink, 64)), text, container.NewCenter(hint)))),
	}
}

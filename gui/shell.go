package main

import (
	"fmt"
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
// its title and subtitle, its icon in the menu, the buttons shown next to
// the title and its content.
type page struct {
	section  string
	title    string
	subtitle string
	icon     fyne.Resource
	actions  []fyne.CanvasObject
	content  fyne.CanvasObject
}

// pageSlide is how far, in pixels, a new page starts from its place before
// it slides in.
const pageSlide = 24

// newShell builds the window content for app a: the sidebar with one item
// per page on the left and, on the right, the title of the current page,
// status chips and the content of the page. Notifications float over it.
func newShell(a *App, pages []page) fyne.CanvasObject {
	title := newHeading("", theme.SizeNameHeadingText)
	subtitle := newCaption("")
	subtitle.TextSize = 13
	actions := container.NewHBox()
	body := newSlider()
	items := make([]*sidebarItem, len(pages))
	current := -1

	// show puts page i on screen and highlights its item in the menu. The
	// page slides in from below when the user goes down the menu and from
	// above when they go up, so the motion points where they went.
	show := func(i int) {
		p := pages[i]
		title.SetText(p.title)
		subtitle.Text = p.subtitle
		subtitle.Refresh()
		actions.Objects = p.actions
		actions.Refresh()
		body.box.Objects = []fyne.CanvasObject{p.content}
		if current >= 0 && i != current {
			from := float32(pageSlide)
			if i < current {
				from = -pageSlide
			}
			body.jump(from)
			body.slideTo(0, nil)
		} else {
			body.box.Refresh()
		}
		current = i
		for j, item := range items {
			item.setSelected(j == i)
		}
	}

	menu := container.NewVBox()
	for i, p := range pages {
		if i == 0 || p.section != pages[i-1].section {
			menu.Add(newSectionTitle(p.section))
		}
		items[i] = newSidebarItem(p.title, p.icon, func() { show(i) })
		menu.Add(items[i])
	}
	show(0)

	heading := container.NewBorder(nil, nil, nil, container.NewCenter(actions),
		container.NewVBox(title, container.New(layout.NewCustomPaddedLayout(0, 0, 8, 0), subtitle)))
	chips := container.New(layout.NewCustomPaddedLayout(10, 14, 4, 0), newStatusChips(a))
	// The header is opaque and is drawn after the body, so a page sliding
	// up goes under it, as content goes under the toolbar of a Mac app.
	header := container.NewStack(canvas.NewRectangle(canvasGray), container.NewVBox(heading, chips))

	// Scroll edge effect: instead of a hard line under the header, the
	// content fades out where it meets it.
	fade := canvas.NewVerticalGradient(canvasGray, tint(canvasGray, 0))
	fade.SetMinSize(fyne.NewSize(0, 14))
	content := container.NewStack(body.box, container.NewBorder(fade, nil, nil, nil))

	right := container.New(layout.NewCustomPaddedLayout(20, 20, 28, 28),
		container.NewBorder(header, nil, nil, nil, content))

	toasts := newToaster()
	a.toasts = toasts
	return container.NewStack(container.NewBorder(nil, nil, newSidebar(menu), nil, right), toasts.layer)
}

// newSidebar returns the gray column on the left: the app icon and name on
// top, the menu in the middle and, at the bottom, the flag of Junior and
// the flag of Colombia.
func newSidebar(menu fyne.CanvasObject) fyne.CanvasObject {
	name := canvas.NewText("Costeños", ink)
	name.TextSize = 16
	name.TextStyle = fyne.TextStyle{Bold: true}
	brand := container.NewBorder(nil, nil, newAppIcon(40), nil,
		container.New(layout.NewCustomPaddedLayout(2, 0, 10, 0), container.NewVBox(name, newCaption("Hospital · Narcolepsia"))))

	cheer := canvas.NewText("¡Junior, tu papá!", juniorRed)
	cheer.TextStyle = fyne.TextStyle{Bold: true}
	pennant := container.NewBorder(nil, nil, container.NewCenter(newJuniorFlag(54, 30)), nil,
		container.New(layout.NewCustomPaddedLayout(0, 0, 10, 0),
			container.NewVBox(layout.NewSpacer(), cheer, newCaption("Rojiblanco de corazón"), layout.NewSpacer())))
	madeIn := container.NewHBox(container.NewCenter(newColombiaFlag(20, 13)), newCaption("Hecho en la Costa"))
	footer := newPanel(container.NewVBox(pennant, widget.NewSeparator(), madeIn))

	content := container.NewBorder(padded(4, brand), footer, nil, nil, container.NewVScroll(menu))

	bg := canvas.NewRectangle(sidebarBg)
	width := canvas.NewRectangle(color.Transparent)
	width.SetMinSize(fyne.NewSize(236, 0)) // fixed width of the sidebar
	line := canvas.NewRectangle(hairline)
	line.SetMinSize(fyne.NewSize(1, 0))
	return container.NewBorder(nil, nil, nil, line, container.NewStack(bg, width, padded(14, content)))
}

// sidebarItem is one entry of the menu: a gray icon and a title. It is a
// small widget of its own because it reacts to the mouse: it darkens as
// soon as the button goes down, runs onTap when it is released and shades
// its background while the mouse is over it.
type sidebarItem struct {
	widget.BaseWidget
	bg       *canvas.Rectangle
	label    *canvas.Text
	icon     *widget.Icon
	normal   fyne.Resource // the icon in gray
	chosen   fyne.Resource // the icon in the accent color, when selected
	onTap    func()
	selected bool
	hovered  bool
	pressed  bool
}

// Compile-time checks: the build fails if sidebarItem stops reacting to
// taps, to the mouse buttons, to hovering or stops showing the hand cursor.
var (
	_ fyne.Tappable      = (*sidebarItem)(nil)
	_ desktop.Mouseable  = (*sidebarItem)(nil)
	_ desktop.Hoverable  = (*sidebarItem)(nil)
	_ desktop.Cursorable = (*sidebarItem)(nil)
)

// newSidebarItem returns a menu item with title and icon that runs onTap
// when tapped.
func newSidebarItem(title string, icon fyne.Resource, onTap func()) *sidebarItem {
	s := &sidebarItem{
		bg:     canvas.NewRectangle(color.Transparent),
		label:  canvas.NewText(title, ink),
		normal: theme.NewColoredResource(icon, colorNameIcon),
		chosen: theme.NewColoredResource(icon, theme.ColorNamePrimary),
		onTap:  onTap,
	}
	s.icon = widget.NewIcon(s.normal)
	s.bg.CornerRadius = 6
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer lays out the item: the icon on the left and the title.
func (s *sidebarItem) CreateRenderer() fyne.WidgetRenderer {
	icon := container.NewCenter(container.NewGridWrap(fyne.NewSize(18, 18), s.icon))
	row := container.NewBorder(nil, nil, icon, nil,
		container.New(layout.NewCustomPaddedLayout(0, 0, 10, 0), container.NewVBox(layout.NewSpacer(), s.label, layout.NewSpacer())))
	return widget.NewSimpleRenderer(container.NewStack(s.bg,
		container.New(layout.NewCustomPaddedLayout(6, 6, 8, 8), row)))
}

// paint sets the background, text and icon of the item from its state: a
// gray piece and a blue icon when selected, as the sidebar of the Finder
// does, a darker shade while pressed and a light one while hovered.
func (s *sidebarItem) paint() {
	s.bg.FillColor = color.Transparent
	s.label.TextStyle = fyne.TextStyle{}
	s.icon.SetResource(s.normal)
	switch {
	case s.selected:
		s.bg.FillColor = selectGray
		s.label.TextStyle = fyne.TextStyle{Bold: true}
		s.icon.SetResource(s.chosen)
	case s.pressed:
		s.bg.FillColor = color.NRGBA{A: 0x1C}
	case s.hovered:
		s.bg.FillColor = color.NRGBA{A: 0x0D}
	}
	s.Refresh()
}

// setSelected marks the item as the current page when on is true.
func (s *sidebarItem) setSelected(on bool) {
	s.selected = on
	s.paint()
}

// Tapped runs the action of the item.
func (s *sidebarItem) Tapped(*fyne.PointEvent) {
	s.onTap()
}

// MouseDown darkens the item at once, before the button is released, so
// the app answers the very moment the user acts.
func (s *sidebarItem) MouseDown(*desktop.MouseEvent) {
	s.pressed = true
	s.paint()
}

// MouseUp removes the pressed shade.
func (s *sidebarItem) MouseUp(*desktop.MouseEvent) {
	s.pressed = false
	s.paint()
}

// MouseIn shades the item while the mouse is over it.
func (s *sidebarItem) MouseIn(*desktop.MouseEvent) {
	s.hovered = true
	s.paint()
}

// MouseMoved does nothing; it is part of desktop.Hoverable.
func (s *sidebarItem) MouseMoved(*desktop.MouseEvent) {}

// MouseOut removes the shade when the mouse leaves the item.
func (s *sidebarItem) MouseOut() {
	s.hovered = false
	s.pressed = false
	s.paint()
}

// Cursor shows a hand over the item, as over a link.
func (s *sidebarItem) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

// counts are the numbers of the hospital shown in the header and the
// summary.
type counts struct {
	patients, inHallway, inBed, freeRooms int
}

// countHospital calculates counts from the getters of h.
func countHospital(h *hospital.Hospital) counts {
	c := counts{patients: len(h.Patients()), inHallway: len(h.PatientsInHallway())}
	for _, p := range h.Patients() {
		if p.State() == hospital.AsleepInBed {
			c.inBed++
		}
	}
	for _, r := range h.Rooms() {
		if r.IsAvailable() {
			c.freeRooms++
		}
	}
	return c
}

// newStatusChips builds the small chips under the title of every page:
// patients, patients asleep in the hallway, patients asleep in bed and free
// rooms. They are the counters of the header, kept small so the page
// itself stays the most important thing on screen.
func newStatusChips(a *App) fyne.CanvasObject {
	patients, inHallway, inBed, freeRooms := newPill(), newPill(), newPill(), newPill()
	a.addRefresher(func() {
		c := countHospital(a.h)
		patients.set(plural(c.patients, "paciente", "pacientes"), accent)
		inHallway.set(fmt.Sprintf("%d en pasillo", c.inHallway), orange)
		inBed.set(fmt.Sprintf("%d en cama", c.inBed), indigo)
		freeRooms.set(plural(c.freeRooms, "cuarto libre", "cuartos libres"), green)
	})
	// Each chip gets some room on its right, so the dots do not crowd.
	spaced := container.NewHBox()
	for _, chip := range []*pill{patients, inHallway, inBed, freeRooms} {
		spaced.Add(container.New(layout.NewCustomPaddedLayout(0, 0, 0, 12), chip.box))
	}
	return spaced
}

// plural returns n followed by one when n is 1 and by many otherwise, as in
// "1 paciente" and "3 pacientes".
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return fmt.Sprintf("%d %s", n, many)
}

// newStatsStrip builds the row of big stat cards of the summary page. It
// registers a refresher that recalculates them from the hospital getters.
func newStatsStrip(a *App) fyne.CanvasObject {
	patients := newStatCard("Pacientes", accent)
	inHallway := newStatCard("Dormidos en pasillo", orange)
	inBed := newStatCard("Dormidos en cama", indigo)
	freeRooms := newStatCard("Cuartos libres", green)

	a.addRefresher(func() {
		c := countHospital(a.h)
		patients.set(c.patients)
		inHallway.set(c.inHallway)
		inBed.set(c.inBed)
		freeRooms.set(c.freeRooms)
	})
	return container.NewGridWithColumns(4, patients.box, inHallway.box, inBed.box, freeRooms.box)
}

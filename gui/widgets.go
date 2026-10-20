package main

import (
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// cardRadius is the corner radius of cards and dialogs.
const cardRadius = 16

// newPanel returns content over a white rounded card with a very soft
// shadow under it, the borderless card look of Apple interfaces.
func newPanel(content fyne.CanvasObject) fyne.CanvasObject {
	shadow := canvas.NewRectangle(color.NRGBA{A: 0x0F})
	shadow.CornerRadius = cardRadius
	bg := canvas.NewRectangle(white)
	bg.CornerRadius = cardRadius
	// The shadow is the same card moved 2 pixels down, so only its bottom
	// edge shows.
	moved := container.New(layout.NewCustomPaddedLayout(2, -2, 0, 0), shadow)
	return container.NewStack(moved, bg, padded(16, content))
}

// newTitledPanel returns a card with a bold title over content.
func newTitledPanel(title string, content fyne.CanvasObject) fyne.CanvasObject {
	return newPanel(container.NewBorder(padded(2, newHeading(title, theme.SizeNameSubHeadingText)), nil, nil, nil, content))
}

// padded returns content with the same padding p on its four sides.
func padded(p float32, content fyne.CanvasObject) fyne.CanvasObject {
	return container.New(layout.NewCustomPaddedLayout(p, p, p, p), content)
}

// newHeading returns a bold label whose text size is the theme size called
// size, such as theme.SizeNameHeadingText.
func newHeading(text string, size fyne.ThemeSizeName) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	l.SizeName = size
	return l
}

// newCaption returns a small gray text for secondary information.
func newCaption(text string) *canvas.Text {
	t := canvas.NewText(text, muted)
	t.TextSize = theme.CaptionTextSize()
	return t
}

// newSectionTitle returns a small bold gray title in capitals, used to
// group items as in the sidebar of macOS.
func newSectionTitle(text string) fyne.CanvasObject {
	t := canvas.NewText(strings.ToUpper(text), faint)
	t.TextSize = 11
	t.TextStyle = fyne.TextStyle{Bold: true}
	return container.New(layout.NewCustomPaddedLayout(10, 2, 8, 0), t)
}

// pill is a small capsule with a soft tint of a color and its text in a
// darker shade of the same color, used for states.
type pill struct {
	bg   *canvas.Rectangle
	text *canvas.Text
	box  *fyne.Container // what is added to a layout
}

// newPill returns an empty pill, ready to be filled with set.
func newPill() *pill {
	p := &pill{
		bg:   canvas.NewRectangle(color.Transparent),
		text: canvas.NewText("", ink),
	}
	p.bg.CornerRadius = 11
	p.text.TextSize = theme.CaptionTextSize()
	p.text.TextStyle = fyne.TextStyle{Bold: true}
	inner := container.New(layout.NewCustomPaddedLayout(3, 3, 10, 10), p.text)
	// Center keeps the pill at its own height inside rows that are taller.
	p.box = container.NewCenter(container.NewStack(p.bg, inner))
	return p
}

// set shows text with the colors of c. Refreshing the box also recalculates
// the layout, because the pill grows or shrinks with its text.
func (p *pill) set(text string, c color.Color) {
	p.text.Text = text
	p.text.Color = darker(c)
	p.bg.FillColor = tint(c, 0x26)
	p.box.Refresh()
}

// avatar is a circle with the initials of a person, in a soft tint of a
// color chosen from their ID.
type avatar struct {
	circle *canvas.Circle
	text   *canvas.Text
	box    *fyne.Container // what is added to a layout
}

// newAvatar returns an empty avatar of size x size pixels.
func newAvatar(size float32) *avatar {
	a := &avatar{
		circle: canvas.NewCircle(fillGray),
		text:   canvas.NewText("", ink),
	}
	a.text.TextSize = size * 0.36
	a.text.TextStyle = fyne.TextStyle{Bold: true}
	circle := container.NewStack(a.circle, container.NewCenter(a.text))
	a.box = container.NewCenter(container.NewGridWrap(fyne.NewSize(size, size), circle))
	return a
}

// set shows the initials of name in the color chosen from id, so the same
// person always gets the same color.
func (a *avatar) set(name, id string) {
	c := avatarColor(id)
	a.text.Text = initials(name)
	a.text.Color = darker(c)
	a.circle.FillColor = tint(c, 0x2E)
	a.box.Refresh()
}

// avatarColors are the colors an avatar can get.
var avatarColors = []color.Color{caribe, orange, indigo, green, pink, sky}

// avatarColor picks one of avatarColors from id: it adds up its characters
// and keeps the remainder, so the result is always the same for one id.
func avatarColor(id string) color.Color {
	sum := 0
	for _, r := range id {
		sum += int(r)
	}
	return avatarColors[sum%len(avatarColors)]
}

// newIconBadge returns icon in white over a rounded square of color c, like
// the icons of the Settings app of Apple.
func newIconBadge(icon fyne.Resource, c color.Color, size float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(c)
	bg.CornerRadius = size * 0.26
	ic := widget.NewIcon(theme.NewColoredResource(icon, theme.ColorNameForegroundOnPrimary))
	return container.NewGridWrap(fyne.NewSize(size, size), container.NewStack(bg, padded(size*0.12, ic)))
}

// statCard is a card with a colored icon, a big number and a caption, in
// the style of the widgets of Apple.
type statCard struct {
	value *canvas.Text
	box   fyne.CanvasObject // what is added to a layout
}

// newStatCard returns a stat card with caption and icon over color c.
func newStatCard(caption string, icon fyne.Resource, c color.Color) *statCard {
	value := canvas.NewText("0", ink)
	value.TextSize = 26
	value.TextStyle = fyne.TextStyle{Bold: true}

	text := container.New(layout.NewCustomPaddedLayout(0, 0, 12, 0), container.NewVBox(value, newCaption(caption)))
	box := newPanel(container.NewBorder(nil, nil, container.NewCenter(newIconBadge(icon, c, 38)), nil, text))
	return &statCard{value: value, box: box}
}

// set shows n as the value of the card.
func (s *statCard) set(n int) {
	s.value.Text = strconv.Itoa(n)
	s.value.Refresh()
}

// newInfoRow returns one row of a grouped list, as in the Settings app: the
// caption in gray on the left and value on the right.
func newInfoRow(caption string, value *widget.Label) fyne.CanvasObject {
	value.Alignment = fyne.TextAlignTrailing
	value.TextStyle = fyne.TextStyle{Bold: true}
	label := widget.NewLabel(caption)
	label.Importance = widget.LowImportance
	return container.NewBorder(nil, nil, label, nil, value)
}

// newGroupedList returns rows inside a gray rounded box, with a hairline
// between them, like an inset grouped list of iOS.
func newGroupedList(rows ...fyne.CanvasObject) fyne.CanvasObject {
	list := container.NewVBox()
	for i, r := range rows {
		if i > 0 {
			list.Add(widget.NewSeparator())
		}
		list.Add(r)
	}
	bg := canvas.NewRectangle(canvasGray)
	bg.CornerRadius = 12
	return container.NewStack(bg, container.New(layout.NewCustomPaddedLayout(4, 4, 10, 10), list))
}

// newSegmented returns a segmented control: one button per label inside a
// gray rounded track. Tapping a segment highlights it and calls onSelect
// with its index. The first segment starts selected.
func newSegmented(labels []string, onSelect func(i int)) fyne.CanvasObject {
	buttons := make([]*widget.Button, len(labels))
	choose := func(i int) {
		for j, b := range buttons {
			b.Importance = widget.LowImportance
			if j == i {
				b.Importance = widget.HighImportance
			}
			b.Refresh()
		}
		onSelect(i)
	}

	row := container.NewGridWithColumns(len(labels))
	for i, l := range labels {
		buttons[i] = widget.NewButton(l, func() { choose(i) })
		row.Add(buttons[i])
	}
	choose(0)

	track := canvas.NewRectangle(fillGray)
	track.CornerRadius = 12
	return container.NewStack(track, padded(3, row))
}

// initials returns the first letter of the first two words of name in upper
// case, for example "CP" for "Carlos Polo".
func initials(name string) string {
	words := strings.Fields(name)
	result := ""
	for i := 0; i < len(words) && i < 2; i++ {
		first := []rune(words[i])[0] // Fields never returns empty words
		result += string(unicode.ToUpper(first))
	}
	return result
}

// setEnabled enables button b when enabled is true and disables it otherwise.
func setEnabled(b *widget.Button, enabled bool) {
	if enabled {
		b.Enable()
	} else {
		b.Disable()
	}
}

// newActionButton returns a button with icon that runs tapped, styled with
// importance, such as widget.HighImportance for the main action of a page.
func newActionButton(label string, icon fyne.Resource, importance widget.Importance, tapped func()) *widget.Button {
	b := widget.NewButtonWithIcon(label, icon, tapped)
	b.Importance = importance
	return b
}

// meter is a thin capsule bar, like the ones of macOS, filled from the left
// up to a part between 0 and 1. It is a widget with its own renderer
// because the width of the fill depends on the size the meter gets. When
// the part changes, the fill moves there with a spring that starts from the
// part on screen (shown), so a new change in the middle never makes it jump.
type meter struct {
	widget.BaseWidget
	shown float32 // part drawn right now, from 0 to 1
	color color.Color
	anim  *fyne.Animation
}

// newMeter returns a meter filled up to value out of total with color c,
// without animation.
func newMeter(value, total float64, c color.Color) *meter {
	m := &meter{shown: part(value, total), color: c}
	m.ExtendBaseWidget(m)
	return m
}

// set moves the fill to value out of total with a spring and paints it with
// color c.
func (m *meter) set(value, total float64, c color.Color) {
	m.color = c
	stopAnimation(m.anim)
	m.anim = animateSpring(m.shown, part(value, total), responseMeter, func(v float32) {
		m.shown = v
		m.Refresh()
	}, nil)
}

// part returns value/total between 0 and 1, or 0 if total is not positive.
func part(value, total float64) float32 {
	if total <= 0 {
		return 0
	}
	return float32(min(max(value/total, 0), 1))
}

// CreateRenderer returns the renderer that draws the meter.
func (m *meter) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(fillGray)
	track.CornerRadius = 4
	fill := canvas.NewRectangle(m.color)
	fill.CornerRadius = 4
	return &meterRenderer{m: m, track: track, fill: fill}
}

// meterRenderer draws a meter: a gray track and a colored fill over it.
type meterRenderer struct {
	m     *meter
	track *canvas.Rectangle
	fill  *canvas.Rectangle
}

// Layout makes the track as big as the meter and the fill as wide as the
// part shown.
func (r *meterRenderer) Layout(size fyne.Size) {
	r.track.Resize(size)
	r.fill.Resize(fyne.NewSize(size.Width*r.m.shown, size.Height))
}

// MinSize is the smallest size of a meter: thin, as on macOS.
func (r *meterRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 8)
}

// Refresh repaints the fill with the current color and width.
func (r *meterRenderer) Refresh() {
	r.fill.FillColor = r.m.color
	r.Layout(r.m.Size())
	r.track.Refresh()
	r.fill.Refresh()
}

// Objects returns the shapes of the meter, from back to front.
func (r *meterRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.track, r.fill}
}

// Destroy does nothing: the meter holds no resources.
func (r *meterRenderer) Destroy() {}

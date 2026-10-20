package main

import (
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// cardRadius is the corner radius of cards and dialogs.
const cardRadius = 12

// newPanel returns content over a white rounded card with a thin gray
// border. Cards are separated from the gray window by their color and that
// hairline, not by colored shadows.
func newPanel(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(white)
	bg.CornerRadius = cardRadius
	bg.StrokeColor = hairline
	bg.StrokeWidth = 1
	return container.NewStack(bg, padded(16, content))
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

// newSectionTitle returns a small bold gray title, used to group items as in
// the sidebar of macOS.
func newSectionTitle(text string) fyne.CanvasObject {
	t := canvas.NewText(text, faint)
	t.TextSize = 11
	t.TextStyle = fyne.TextStyle{Bold: true}
	return container.New(layout.NewCustomPaddedLayout(12, 4, 10, 0), t)
}

// pill shows a state as macOS Mail or Reminders do: a small colored dot and
// gray text. The color lives only in the dot, so a list full of states
// stays calm.
type pill struct {
	dot  *canvas.Circle
	text *canvas.Text
	box  *fyne.Container // what is added to a layout
}

// newPill returns an empty pill, ready to be filled with set.
func newPill() *pill {
	p := &pill{
		dot:  canvas.NewCircle(gray),
		text: canvas.NewText("", muted),
	}
	p.text.TextSize = theme.CaptionTextSize()
	dot := container.NewCenter(container.NewGridWrap(fyne.NewSize(7, 7), p.dot))
	// Center keeps the pill at its own height inside rows that are taller.
	p.box = container.NewCenter(container.NewHBox(dot, p.text))
	return p
}

// set shows text with a dot of color c. Refreshing the box also
// recalculates the layout, because the pill grows or shrinks with its text.
func (p *pill) set(text string, c color.Color) {
	p.text.Text = text
	p.dot.FillColor = c
	p.box.Refresh()
}

// avatar is a gray circle with the initials of a person in white, like the
// monograms of the Contacts app.
type avatar struct {
	text *canvas.Text
	box  *fyne.Container // what is added to a layout
}

// newAvatar returns an empty avatar of size x size pixels.
func newAvatar(size float32) *avatar {
	a := &avatar{text: canvas.NewText("", white)}
	a.text.TextSize = size * 0.38
	a.text.TextStyle = fyne.TextStyle{Bold: true}
	circle := container.NewStack(canvas.NewCircle(monogram), container.NewCenter(a.text))
	a.box = container.NewCenter(container.NewGridWrap(fyne.NewSize(size, size), circle))
	return a
}

// set shows the initials of name.
func (a *avatar) set(name string) {
	a.text.Text = initials(name)
	a.box.Refresh()
}

// statCard is a card with a caption over a big number.
type statCard struct {
	value *canvas.Text
	box   fyne.CanvasObject // what is added to a layout
}

// newStatCard returns a stat card with caption and a dot of color c, the
// only color of the card.
func newStatCard(caption string, c color.Color) *statCard {
	value := canvas.NewText("0", ink)
	value.TextSize = 28
	value.TextStyle = fyne.TextStyle{Bold: true}
	label := newPill()
	label.set(caption, c)
	box := newPanel(container.NewVBox(container.NewHBox(label.box), value))
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
	label := widget.NewLabel(caption)
	label.Importance = widget.LowImportance
	return container.NewBorder(nil, nil, label, nil, value)
}

// newGroupedList returns rows inside a light gray rounded box, with a
// hairline between them, like an inset grouped list of iOS.
func newGroupedList(rows ...fyne.CanvasObject) fyne.CanvasObject {
	list := container.NewVBox()
	for i, r := range rows {
		if i > 0 {
			list.Add(widget.NewSeparator())
		}
		list.Add(r)
	}
	bg := canvas.NewRectangle(fieldGray)
	bg.CornerRadius = 10
	return container.NewStack(bg, container.New(layout.NewCustomPaddedLayout(2, 2, 10, 10), list))
}

// newSegmented returns a segmented control as macOS draws it: the options
// over a gray track, and the chosen one on a white rounded piece. Tapping an
// option chooses it and calls onSelect with its index. The first option
// starts chosen.
func newSegmented(labels []string, onSelect func(i int)) fyne.CanvasObject {
	segments := make([]*segment, len(labels))
	choose := func(i int) {
		for j, s := range segments {
			s.setSelected(j == i)
		}
		onSelect(i)
	}

	row := container.NewGridWithColumns(len(labels))
	for i, l := range labels {
		segments[i] = newSegment(l, func() { choose(i) })
		row.Add(segments[i])
	}
	choose(0)

	track := canvas.NewRectangle(fillGray)
	track.CornerRadius = 9
	return container.NewStack(track, padded(2, row))
}

// segment is one option of a segmented control. It is a small widget so it
// can be tapped and show the hand cursor.
type segment struct {
	widget.BaseWidget
	bg    *canvas.Rectangle
	label *canvas.Text
	onTap func()
}

// Compile-time checks: the build fails if segment stops reacting to taps or
// stops showing the hand cursor.
var (
	_ fyne.Tappable      = (*segment)(nil)
	_ desktop.Cursorable = (*segment)(nil)
)

// newSegment returns an option with title that runs onTap when tapped.
func newSegment(title string, onTap func()) *segment {
	s := &segment{bg: canvas.NewRectangle(color.Transparent), label: canvas.NewText(title, muted), onTap: onTap}
	s.bg.CornerRadius = 7
	s.label.TextSize = 13
	s.ExtendBaseWidget(s)
	return s
}

// CreateRenderer lays out the option: its title centered over the piece.
func (s *segment) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(s.bg,
		container.New(layout.NewCustomPaddedLayout(6, 6, 8, 8), container.NewCenter(s.label))))
}

// setSelected puts the option on a white piece with a hairline border when
// on is true.
func (s *segment) setSelected(on bool) {
	s.bg.FillColor = color.Transparent
	s.bg.StrokeWidth = 0
	s.label.Color = muted
	s.label.TextStyle = fyne.TextStyle{}
	if on {
		s.bg.FillColor = white
		s.bg.StrokeColor = hairline
		s.bg.StrokeWidth = 1
		s.label.Color = ink
		s.label.TextStyle = fyne.TextStyle{Bold: true}
	}
	s.Refresh()
}

// Tapped chooses the option.
func (s *segment) Tapped(*fyne.PointEvent) {
	s.onTap()
}

// Cursor shows a hand over the option.
func (s *segment) Cursor() desktop.Cursor {
	return desktop.PointerCursor
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

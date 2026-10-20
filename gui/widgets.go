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

// panelRadius is the corner radius of panels, cards and dialogs.
const panelRadius = 14

// newPanel returns content over a white rounded rectangle with a thin
// border: the card look used by every page.
func newPanel(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(white)
	bg.CornerRadius = panelRadius
	bg.StrokeColor = border
	bg.StrokeWidth = 1
	return container.NewStack(bg, padded(14, content))
}

// newTitledPanel returns a panel with a bold title over content.
func newTitledPanel(title string, content fyne.CanvasObject) fyne.CanvasObject {
	return newPanel(container.NewBorder(newHeading(title, theme.SizeNameSubHeadingText), nil, nil, nil, content))
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

// pill is a small rounded label with a background color, used for states.
type pill struct {
	bg   *canvas.Rectangle
	text *canvas.Text
	box  *fyne.Container // what is added to a layout
}

// newPill returns an empty pill, ready to be filled with set.
func newPill() *pill {
	p := &pill{
		bg:   canvas.NewRectangle(color.Transparent),
		text: canvas.NewText("", white),
	}
	p.bg.CornerRadius = 10
	p.text.TextSize = theme.CaptionTextSize()
	p.text.TextStyle = fyne.TextStyle{Bold: true}
	inner := container.New(layout.NewCustomPaddedLayout(3, 3, 10, 10), p.text)
	// Center keeps the pill at its own height inside rows that are taller.
	p.box = container.NewCenter(container.NewStack(p.bg, inner))
	return p
}

// set shows text over color c. Refreshing the box also recalculates the
// layout, because the pill grows or shrinks with its text.
func (p *pill) set(text string, c color.Color) {
	p.text.Text = text
	p.bg.FillColor = c
	p.box.Refresh()
}

// avatar is a colored circle with the initials of a person.
type avatar struct {
	circle *canvas.Circle
	text   *canvas.Text
	box    *fyne.Container // what is added to a layout
}

// newAvatar returns an empty avatar of size x size pixels.
func newAvatar(size float32) *avatar {
	a := &avatar{
		circle: canvas.NewCircle(turquoise),
		text:   canvas.NewText("", white),
	}
	a.text.TextSize = size * 0.38
	a.text.TextStyle = fyne.TextStyle{Bold: true}
	circle := container.NewStack(a.circle, container.NewCenter(a.text))
	a.box = container.NewCenter(container.NewGridWrap(fyne.NewSize(size, size), circle))
	return a
}

// set shows the initials of name over a color chosen from id, so the same
// person always gets the same color.
func (a *avatar) set(name, id string) {
	a.text.Text = initials(name)
	a.circle.FillColor = avatarColor(id)
	a.box.Refresh()
}

// avatarColors are the colors an avatar can get.
var avatarColors = []color.Color{turquoise, coral, indigo, amber, deepSea, green}

// avatarColor picks one of avatarColors from id: it adds up its characters
// and keeps the remainder, so the result is always the same for one id.
func avatarColor(id string) color.Color {
	sum := 0
	for _, r := range id {
		sum += int(r)
	}
	return avatarColors[sum%len(avatarColors)]
}

// statCard is a panel with a colored icon, a big number and a caption.
type statCard struct {
	value *canvas.Text
	box   fyne.CanvasObject // what is added to a layout
}

// newStatCard returns a stat card with caption and icon over color c.
func newStatCard(caption string, icon fyne.Resource, c color.Color) *statCard {
	circle := canvas.NewCircle(c)
	ic := widget.NewIcon(theme.NewColoredResource(icon, theme.ColorNameForegroundOnPrimary))
	badge := container.NewGridWrap(fyne.NewSize(44, 44), container.NewStack(circle, padded(10, ic)))

	value := canvas.NewText("0", ink)
	value.TextSize = 26
	value.TextStyle = fyne.TextStyle{Bold: true}
	text := container.NewVBox(value, newCaption(caption))

	box := newPanel(container.NewBorder(nil, nil, container.NewCenter(badge), nil, padded(4, text)))
	return &statCard{value: value, box: box}
}

// set shows n as the value of the card.
func (s *statCard) set(n int) {
	s.value.Text = strconv.Itoa(n)
	s.value.Refresh()
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

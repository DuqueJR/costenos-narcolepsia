package main

import (
	"image/color"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"costenos-narcolepsia/hospital"
)

// Palette of the GUI, taken from the light appearance of macOS. It is almost
// all neutral grays with one accent color, the blue of Apple. The other
// colors only appear as small status dots, where they mean something.
var (
	canvasGray  = color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF7, A: 0xFF} // window background
	sidebarBg   = color.NRGBA{R: 0xEE, G: 0xEE, B: 0xF0, A: 0xFF} // sidebar
	selectGray  = color.NRGBA{R: 0xDE, G: 0xDE, B: 0xE3, A: 0xFF} // selected sidebar item
	fillGray    = color.NRGBA{R: 0xE8, G: 0xE8, B: 0xED, A: 0xFF} // buttons and segmented track
	fieldGray   = color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF7, A: 0xFF} // text fields and grouped lists
	hairline    = color.NRGBA{R: 0xE5, G: 0xE5, B: 0xEA, A: 0xFF} // card borders and separators
	fieldBorder = color.NRGBA{R: 0xD2, G: 0xD2, B: 0xD7, A: 0xFF} // borders of fields and checkboxes
	white       = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // cards
	ink         = color.NRGBA{R: 0x1D, G: 0x1D, B: 0x1F, A: 0xFF} // text
	muted       = color.NRGBA{R: 0x6E, G: 0x6E, B: 0x73, A: 0xFF} // secondary text and icons
	faint       = color.NRGBA{R: 0xAE, G: 0xAE, B: 0xB2, A: 0xFF} // placeholders and disabled text
	monogram    = color.NRGBA{R: 0xA8, G: 0xA8, B: 0xAD, A: 0xFF} // avatars, as the monograms of Contacts

	accent = color.NRGBA{R: 0x00, G: 0x71, B: 0xE3, A: 0xFF} // the only accent: the blue of Apple
	cream  = color.NRGBA{R: 0xF4, G: 0xEA, B: 0xD5, A: 0xFF} // cane fiber of the sombrero vueltiao

	green  = color.NRGBA{R: 0x34, G: 0xC7, B: 0x59, A: 0xFF}
	orange = color.NRGBA{R: 0xFF, G: 0x95, B: 0x00, A: 0xFF}
	indigo = color.NRGBA{R: 0x58, G: 0x56, B: 0xD6, A: 0xFF}
	red    = color.NRGBA{R: 0xFF, G: 0x3B, B: 0x30, A: 0xFF}
	gray   = color.NRGBA{R: 0x8E, G: 0x8E, B: 0x93, A: 0xFF}
)

// colorNameIcon is a color name of our own for the icons of the sidebar.
// Fyne asks the theme for it like for any other name, so the icons can be
// recolored with theme.NewColoredResource.
const colorNameIcon fyne.ThemeColorName = "costenos-icon"

// caribeTheme is the light theme of the GUI. It changes some colors and
// sizes and uses the font of the system; icons and every other value come
// from the default Fyne theme. It ignores the variant, so the GUI stays
// light even if the system uses dark mode.
type caribeTheme struct {
	regular, bold, italic, boldItalic fyne.Resource // nil means the Fyne font
}

// newCaribeTheme returns the theme with the font of the system when it can
// find it: Segoe UI, the font of Windows apps, as Apple apps use the font of
// their system. On other systems, or if a file is missing, that style keeps
// the font that comes with Fyne.
func newCaribeTheme() caribeTheme {
	var t caribeTheme
	if runtime.GOOS != "windows" {
		return t
	}
	fonts := filepath.Join(os.Getenv("WINDIR"), "Fonts")
	load := func(name string) fyne.Resource {
		r, err := fyne.LoadResourceFromPath(filepath.Join(fonts, name))
		if err != nil {
			return nil // not a failure: Font falls back to the Fyne font
		}
		return r
	}
	t.regular = load("segoeui.ttf")
	// Semibold instead of bold: Apple titles use a semibold weight, which
	// gives presence without looking heavy.
	t.bold = load("segoeuisb.ttf")
	if t.bold == nil {
		t.bold = load("segoeuib.ttf")
	}
	t.italic = load("segoeuii.ttf")
	t.boldItalic = load("segoeuiz.ttf")
	return t
}

// Compile-time check: the build fails if caribeTheme stops being a fyne.Theme.
var _ fyne.Theme = caribeTheme{}

// Color returns the color of the theme called name.
func (caribeTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground, theme.ColorNameHeaderBackground:
		return canvasGray
	case theme.ColorNameButton, theme.ColorNameDisabledButton:
		return fillGray
	case theme.ColorNameInputBackground:
		return fieldGray
	case theme.ColorNameInputBorder:
		return fieldBorder // a subtle border, as macOS text fields and checkboxes have
	case theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return white
	case theme.ColorNameSeparator:
		return hairline
	case theme.ColorNamePrimary, theme.ColorNameHyperlink:
		return accent
	case theme.ColorNameFocus:
		return tint(accent, 0x80)
	case theme.ColorNameSelection:
		return tint(accent, 0x18)
	case colorNameIcon:
		return muted
	case theme.ColorNameHover:
		return color.NRGBA{A: 0x0A}
	case theme.ColorNamePressed:
		return color.NRGBA{A: 0x16}
	case theme.ColorNameShadow:
		return color.NRGBA{A: 0x22}
	case theme.ColorNameScrollBar:
		return color.NRGBA{A: 0x3C}
	case theme.ColorNameForeground:
		return ink
	case theme.ColorNamePlaceHolder, theme.ColorNameDisabled:
		return faint
	case theme.ColorNameError:
		return red
	case theme.ColorNameWarning:
		return orange
	case theme.ColorNameSuccess:
		return green
	case theme.ColorNameForegroundOnPrimary, theme.ColorNameForegroundOnError:
		return white
	}
	return theme.DefaultTheme().Color(name, theme.VariantLight)
}

// Font returns the system font for style s, or the Fyne font when there is
// none. Monospace and symbol text always use the Fyne fonts.
func (t caribeTheme) Font(s fyne.TextStyle) fyne.Resource {
	var font fyne.Resource
	switch {
	case s.Monospace, s.Symbol:
	case s.Bold && s.Italic:
		font = t.boldItalic
	case s.Bold:
		font = t.bold
	case s.Italic:
		font = t.italic
	default:
		font = t.regular
	}
	if font == nil {
		return theme.DefaultTheme().Font(s)
	}
	return font
}

// Icon returns the default Fyne icon called name.
func (caribeTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size of the theme called name: rounder corners, more
// space inside widgets and large titles, as in Apple interfaces.
func (caribeTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius, theme.SizeNameButtonRadius:
		return 8
	case theme.SizeNameCardRadius, theme.SizeNameDialogRadius, theme.SizeNamePopupRadius:
		return cardRadius
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNamePadding:
		return 6 // a little more air between cards and widgets than Fyne's 4
	case theme.SizeNameText:
		return 14
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameHeadingText:
		return 30
	case theme.SizeNameSubHeadingText:
		return 17
	}
	return theme.DefaultTheme().Size(name)
}

// tint returns the opaque color c with alpha a, for soft backgrounds such
// as the one of a pill.
func tint(c color.Color, a uint8) color.Color {
	r, g, b, _ := c.RGBA() // 16-bit channels; our colors are opaque
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: a}
}

// stateColor returns the color that identifies patient state s in the GUI:
// green when awake, orange (amber) when asleep in the hallway and indigo
// when asleep in bed.
func stateColor(s hospital.PatientState) color.Color {
	switch s {
	case hospital.Awake:
		return green
	case hospital.AsleepInHallway:
		return orange
	case hospital.AsleepInBed:
		return indigo
	default:
		return gray // unknown state
	}
}

// stateLabel returns patient state s in Spanish, for the GUI. The domain
// keeps its own English String, used by the console program.
func stateLabel(s hospital.PatientState) string {
	switch s {
	case hospital.Awake:
		return "Despierto"
	case hospital.AsleepInHallway:
		return "Dormido en pasillo"
	case hospital.AsleepInBed:
		return "Dormido en cama"
	default:
		return s.String()
	}
}

// levelColor returns the color of narcolepsy level l: the more severe, the
// warmer.
func levelColor(l hospital.NarcolepsyLevel) color.Color {
	switch l {
	case hospital.Mild:
		return gray
	case hospital.Moderate:
		return orange
	case hospital.Severe:
		return red
	default:
		return gray
	}
}

// levelLabel returns narcolepsy level l in Spanish, for the GUI.
func levelLabel(l hospital.NarcolepsyLevel) string {
	switch l {
	case hospital.Mild:
		return "Leve"
	case hospital.Moderate:
		return "Moderada"
	case hospital.Severe:
		return "Severa"
	default:
		return l.String()
	}
}

// roomStateColor returns green for an available room and indigo for an
// occupied one.
func roomStateColor(s hospital.RoomState) color.Color {
	switch s {
	case hospital.Available:
		return green
	case hospital.Occupied:
		return indigo
	default:
		return gray // unknown state
	}
}

// roomStateLabel returns room state s in Spanish, for the GUI.
func roomStateLabel(s hospital.RoomState) string {
	switch s {
	case hospital.Available:
		return "Disponible"
	case hospital.Occupied:
		return "Ocupado"
	default:
		return s.String()
	}
}

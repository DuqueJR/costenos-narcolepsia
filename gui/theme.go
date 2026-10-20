package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"costenos-narcolepsia/hospital"
)

// Palette of the GUI: the clean grays of Apple interfaces plus the colors of
// the Caribbean coast, the sea and the sun.
var (
	canvasGray = color.NRGBA{R: 0xF5, G: 0xF5, B: 0xF7, A: 0xFF} // window background
	sidebarBg  = color.NRGBA{R: 0xEE, G: 0xEE, B: 0xF2, A: 0xFF} // sidebar
	fillGray   = color.NRGBA{R: 0xE9, G: 0xE9, B: 0xEC, A: 0xFF} // buttons and segmented track
	fieldGray  = color.NRGBA{R: 0xF2, G: 0xF2, B: 0xF5, A: 0xFF} // text fields
	hairline   = color.NRGBA{R: 0xE3, G: 0xE3, B: 0xE8, A: 0xFF} // separators
	white      = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // cards
	ink        = color.NRGBA{R: 0x1D, G: 0x1D, B: 0x1F, A: 0xFF} // text
	muted      = color.NRGBA{R: 0x6E, G: 0x6E, B: 0x73, A: 0xFF} // secondary text
	faint      = color.NRGBA{R: 0xA1, G: 0xA1, B: 0xA6, A: 0xFF} // placeholders and disabled text

	caribe = color.NRGBA{R: 0x0A, G: 0x7E, B: 0xA4, A: 0xFF} // main color: the Caribbean sea
	sky    = color.NRGBA{R: 0x3F, G: 0xB8, B: 0xD9, A: 0xFF} // light end of the sea gradient
	sun    = color.NRGBA{R: 0xFF, G: 0xC2, B: 0x33, A: 0xFF} // the coast sun
	cream  = color.NRGBA{R: 0xF4, G: 0xEA, B: 0xD5, A: 0xFF} // cane fiber of the sombrero vueltiao

	green  = color.NRGBA{R: 0x34, G: 0xC7, B: 0x59, A: 0xFF}
	orange = color.NRGBA{R: 0xFF, G: 0x95, B: 0x00, A: 0xFF}
	indigo = color.NRGBA{R: 0x58, G: 0x56, B: 0xD6, A: 0xFF}
	red    = color.NRGBA{R: 0xFF, G: 0x45, B: 0x3A, A: 0xFF}
	pink   = color.NRGBA{R: 0xFF, G: 0x2D, B: 0x55, A: 0xFF}
	gray   = color.NRGBA{R: 0x8E, G: 0x8E, B: 0x93, A: 0xFF}
)

// caribeTheme is the light theme of the GUI. It only changes some colors
// and sizes; fonts, icons and every other value come from the default Fyne
// theme. It ignores the variant, so the GUI stays light even if the system
// uses dark mode.
type caribeTheme struct{}

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
		return fieldGray // borderless fields; the focus color still draws a ring
	case theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return white
	case theme.ColorNameSeparator:
		return hairline
	case theme.ColorNamePrimary, theme.ColorNameHyperlink:
		return caribe
	case theme.ColorNameFocus:
		return tint(caribe, 0x80)
	case theme.ColorNameSelection:
		return tint(caribe, 0x1F)
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

// Font returns the default Fyne font for style s.
func (caribeTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

// Icon returns the default Fyne icon called name.
func (caribeTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size of the theme called name: rounder corners, more
// space inside widgets and large titles, as in Apple interfaces.
func (caribeTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius:
		return 8
	case theme.SizeNameButtonRadius:
		return 10
	case theme.SizeNameCardRadius, theme.SizeNameDialogRadius, theme.SizeNamePopupRadius:
		return cardRadius
	case theme.SizeNameInnerPadding:
		return 10
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

// darker returns the opaque color c a quarter darker, so text over a tint
// of c is still easy to read.
func darker(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8((r >> 8) * 3 / 4), G: uint8((g >> 8) * 3 / 4), B: uint8((b >> 8) * 3 / 4), A: 0xFF}
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
		return caribe
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

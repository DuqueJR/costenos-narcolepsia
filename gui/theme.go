package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"costenos-narcolepsia/hospital"
)

// Palette of the GUI: a light Caribbean look with sand, turquoise and coral.
var (
	sand      = color.NRGBA{R: 0xFA, G: 0xF6, B: 0xEF, A: 0xFF} // window background
	sandDark  = color.NRGBA{R: 0xF1, G: 0xEB, B: 0xE1, A: 0xFF} // buttons and table headers
	white     = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // panels
	border    = color.NRGBA{R: 0xE8, G: 0xE1, B: 0xD5, A: 0xFF} // panel borders and separators
	turquoise = color.NRGBA{R: 0x0E, G: 0x9F, B: 0x9A, A: 0xFF} // main color
	deepSea   = color.NRGBA{R: 0x0B, G: 0x4F, B: 0x5C, A: 0xFF} // dark end of the banner
	coral     = color.NRGBA{R: 0xF2, G: 0x66, B: 0x5B, A: 0xFF} // errors and dangerous actions
	ink       = color.NRGBA{R: 0x1F, G: 0x2D, B: 0x3D, A: 0xFF} // text
	muted     = color.NRGBA{R: 0x6B, G: 0x7A, B: 0x89, A: 0xFF} // secondary text

	green  = color.NRGBA{R: 0x16, G: 0xA3, B: 0x4A, A: 0xFF}
	amber  = color.NRGBA{R: 0xD9, G: 0x77, B: 0x06, A: 0xFF}
	indigo = color.NRGBA{R: 0x4F, G: 0x46, B: 0xE5, A: 0xFF}
	gray   = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	slate  = color.NRGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF}
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
	case theme.ColorNameBackground:
		return sand
	case theme.ColorNameButton, theme.ColorNameHeaderBackground:
		return sandDark
	case theme.ColorNameInputBackground, theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return white
	case theme.ColorNameInputBorder, theme.ColorNameSeparator:
		return border
	case theme.ColorNamePrimary, theme.ColorNameHyperlink:
		return turquoise
	case theme.ColorNameFocus:
		return color.NRGBA{R: turquoise.R, G: turquoise.G, B: turquoise.B, A: 0x88}
	case theme.ColorNameSelection:
		return color.NRGBA{R: turquoise.R, G: turquoise.G, B: turquoise.B, A: 0x33}
	case theme.ColorNameHover:
		return color.NRGBA{R: turquoise.R, G: turquoise.G, B: turquoise.B, A: 0x1A}
	case theme.ColorNameForeground:
		return ink
	case theme.ColorNamePlaceHolder:
		return muted
	case theme.ColorNameError:
		return coral
	case theme.ColorNameWarning:
		return amber
	case theme.ColorNameSuccess:
		return green
	case theme.ColorNameForegroundOnPrimary, theme.ColorNameForegroundOnError:
		return white
	case theme.ColorNameDisabledButton:
		return sandDark
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0xA9, G: 0xB3, B: 0xBD, A: 0xFF}
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

// Size returns the size of the theme called name: rounder corners and
// bigger headings than the default theme.
func (caribeTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameButtonRadius, theme.SizeNameSelectionRadius:
		return 8
	case theme.SizeNameCardRadius, theme.SizeNameDialogRadius:
		return panelRadius
	case theme.SizeNameHeadingText:
		return 26
	case theme.SizeNameSubHeadingText:
		return 18
	}
	return theme.DefaultTheme().Size(name)
}

// stateColor returns the color that identifies patient state s in the GUI:
// green when awake, amber when asleep in the hallway and indigo when asleep
// in bed.
func stateColor(s hospital.PatientState) color.Color {
	switch s {
	case hospital.Awake:
		return green
	case hospital.AsleepInHallway:
		return amber
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
		return slate
	case hospital.Moderate:
		return amber
	case hospital.Severe:
		return coral
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

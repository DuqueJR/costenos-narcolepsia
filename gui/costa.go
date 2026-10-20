package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// The costeño details of the GUI. The sombrero and the Junior flag are
// pictures (see images/CREDITS.md). The vueltiao stripe and the flag of
// Colombia are canvas.Raster objects: Fyne calls a Go function once per
// pixel and paints the color it returns, so they need no image files.

// noColor is a fully transparent pixel. It must be a color.NRGBA, like the
// other pixels: Fyne picks the kind of image from the type of the first
// pixel, and color.Transparent (a color.Alpha16) would make an image that
// only stores transparency and loses every color.
var noColor = color.NRGBA{}

// Colors of the costeño drawings.
var (
	juniorRed = color.NRGBA{R: 0xD7, G: 0x19, B: 0x2D, A: 0xFF} // the red stripes of Junior de Barranquilla
	flagBlue  = color.NRGBA{R: 0x00, G: 0x38, B: 0x93, A: 0xFF} // blue of the flag of Colombia
	flagRed   = color.NRGBA{R: 0xCE, G: 0x11, B: 0x26, A: 0xFF} // red of the flag of Colombia
	flagGold  = color.NRGBA{R: 0xFC, G: 0xD1, B: 0x16, A: 0xFF} // yellow of the flag of Colombia
)

// newRaster returns a drawing of width x height painted by pixel.
func newRaster(width, height float32, pixel func(x, y, w, h int) color.Color) fyne.CanvasObject {
	r := canvas.NewRasterWithPixels(pixel)
	r.SetMinSize(fyne.NewSize(width, height))
	return r
}

// newVueltiaoBand returns a strip with the pattern of the sombrero vueltiao,
// the woven cane hat of the Colombian Caribbean coast.
func newVueltiaoBand(width, height float32) fyne.CanvasObject {
	return newRaster(width, height, func(x, y, _, h int) color.Color {
		return vueltiaoPixel(x, y, h)
	})
}

// vueltiaoPixel returns the color of pixel (x, y) of the vueltiao pattern:
// cream cane fiber with two dark zigzags that cross and make diamonds. tile
// is the height of one zigzag: the pattern repeats every tile pixels down
// and every 2*tile pixels across.
func vueltiaoPixel(x, y, tile int) color.Color {
	if tile < 2 {
		return cream
	}
	zig := x % (2 * tile)
	if zig > tile {
		zig = 2*tile - zig // goes 0 → tile → 0, a zigzag
	}
	y %= tile
	thick := max(tile/7, 1)
	if absInt(y-zig) <= thick || absInt(y-(tile-zig)) <= thick {
		return ink
	}
	return cream
}

// newSombrero returns the sombrero vueltiao, the woven hat of the coast, in
// a box of width x height. The picture keeps its proportions inside it.
func newSombrero(width, height float32) fyne.CanvasObject {
	img := canvas.NewImageFromResource(picture("sombrero"))
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}

// newAppIcon returns the icon of the app: the sombrero vueltiao on a white
// rounded square with a thin border, like a macOS icon.
func newAppIcon(size float32) fyne.CanvasObject {
	bg := canvas.NewRectangle(white)
	bg.CornerRadius = size * 0.22
	bg.StrokeColor = fieldBorder
	bg.StrokeWidth = 1
	hat := newSombrero(size*0.8, size*0.8)
	return container.NewGridWrap(fyne.NewSize(size, size), container.NewStack(bg, padded(size*0.1, hat)))
}

// newJuniorFlag returns the flag of Junior de Barranquilla, the football
// team of the city, with slightly rounded corners.
func newJuniorFlag(width, height float32) fyne.CanvasObject {
	img := canvas.NewImageFromResource(picture("junior"))
	img.FillMode = canvas.ImageFillContain
	img.CornerRadius = 3
	img.SetMinSize(fyne.NewSize(width, height))
	return img
}

// newColombiaFlag returns the flag of Colombia: yellow on the top half,
// then blue and red.
func newColombiaFlag(width, height float32) fyne.CanvasObject {
	return newRaster(width, height, func(x, y, w, h int) color.Color {
		if !insideRounded(x, y, w, h, max(h/6, 1)) {
			return noColor
		}
		switch fy := float64(y) / float64(h); {
		case fy < 0.5:
			return flagGold
		case fy < 0.75:
			return flagBlue
		default:
			return flagRed
		}
	})
}

// insideRounded reports whether pixel (x, y) is inside a rectangle of
// w x h pixels with corners of radius r: it measures the distance from the
// pixel to the closest point of the inner rectangle that has no corners.
func insideRounded(x, y, w, h, r int) bool {
	cx := min(max(x, r), w-1-r)
	cy := min(max(y, r), h-1-r)
	dx, dy := x-cx, y-cy
	return dx*dx+dy*dy <= r*r
}

// absInt returns the absolute value of n.
func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Names used by spanishDate. The time package only knows English names.
var (
	weekdays = [...]string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
	months   = [...]string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio",
		"agosto", "septiembre", "octubre", "noviembre", "diciembre"}
)

// spanishDate returns t as a date in Spanish, for example
// "martes 20 de octubre".
func spanishDate(t time.Time) string {
	return fmt.Sprintf("%s %d de %s", weekdays[t.Weekday()], t.Day(), months[t.Month()-1])
}

// sayings are costeño sayings for the "phrase of the day" of the dashboard.
var sayings = []string{
	"Quien lo vive es quien lo goza.",
	"¡Ajá, y entonces! ¿Qué más pues?",
	"Más perdío que el hijo de Lindbergh.",
	"Tranquilo, mi llave, que aquí no pasa nada… solo se duermen.",
	"¡Eche! Dormirse en la hamaca sí, en el pasillo no.",
	"Ni siesta que dure cien años, ni paciente que no la aguante.",
	"Erda, qué calor… y el paciente tan fresco, dormido.",
	"¡Junior, tu papá!",
	"Sin afán, que el que se afana se cansa.",
	"Ahí vamos, ahí vamos… pa' lante como el cangrejo, pero derecho.",
}

// sayingOfTheDay returns one of sayings for day t: it changes every day and
// is the same all day long.
func sayingOfTheDay(t time.Time) string {
	return sayings[t.YearDay()%len(sayings)]
}

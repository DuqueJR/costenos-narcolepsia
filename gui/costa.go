package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// The drawings of this file are canvas.Raster objects: Fyne calls a Go
// function once per pixel and paints the color it returns. They are pure
// Go, with no image files. Most functions work with fx and fy, the position
// of the pixel from 0 to 1, so a drawing does not depend on its size.

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
	pole      = color.NRGBA{R: 0x5C, G: 0x4A, B: 0x3A, A: 0xFF} // wooden stick of the pennant
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

// newSombrero returns a sombrero vueltiao of width x height, seen from the
// front: a wide brim with woven rings and a tall woven crown.
func newSombrero(width, height float32) fyne.CanvasObject {
	return newRaster(width, height, sombreroPixel)
}

// sombreroPixel returns the color of pixel (x, y) of a sombrero vueltiao of
// w x h pixels. The brim and the crown are two ellipses: a point is inside
// an ellipse when ((fx-cx)/rx)² + ((fy-cy)/ry)² <= 1. When the hat is
// small the woven pattern would look like a smudge, so it uses stripes.
func sombreroPixel(x, y, w, h int) color.Color {
	fx := float64(x) / float64(w)
	fy := float64(y) / float64(h)
	tile := h / 7
	woven := func() color.Color {
		if tile < 4 {
			if int(fy*14)%2 == 0 {
				return ink
			}
			return cream
		}
		return vueltiaoPixel(x, y, tile)
	}

	cx, cy := (fx-0.5)/0.21, (fy-0.43)/0.36
	if cx*cx+cy*cy <= 1 && fy < 0.72 { // the crown, above the brim
		if fy > 0.6 {
			return ink // the dark band around the base of the crown
		}
		return woven()
	}

	bx, by := (fx-0.5)/0.49, (fy-0.75)/0.2
	d := math.Sqrt(bx*bx + by*by) // 1 on the edge of the brim
	switch {
	case d > 1:
		return noColor
	case d > 0.92:
		return ink // the dark edge of the brim
	case d > 0.6 && d < 0.84:
		return woven() // the woven ring of the brim
	default:
		return cream
	}
}

// newAppIcon returns the icon of the app: a rounded square with the
// Caribbean sea and a sombrero vueltiao in the middle, like a macOS icon.
func newAppIcon(size float32) fyne.CanvasObject {
	return newRaster(size, size, func(x, y, w, h int) color.Color {
		if !insideRounded(x, y, w, h, w*22/100) {
			return noColor
		}
		// The hat is drawn in the box that goes from 10% to 90% across and
		// from 22% to 82% down; outside it, or where the hat is transparent,
		// the sea shows.
		x0, y0, hw, hh := w/10, h*22/100, w*8/10, h*6/10
		if x >= x0 && x < x0+hw && y >= y0 && y < y0+hh {
			if c := sombreroPixel(x-x0, y-y0, hw, hh); c != noColor {
				return c
			}
		}
		return mix(sky, caribe, float64(y)/float64(h))
	})
}

// newJuniorPennant returns a pennant with the red and white stripes of
// Junior de Barranquilla, the football team of the city, on a stick.
func newJuniorPennant(width, height float32) fyne.CanvasObject {
	return newRaster(width, height, func(x, y, w, h int) color.Color {
		stick := max(w/16, 2)
		if x < stick {
			return pole
		}
		// A triangle that points right: it gets thinner as fx grows.
		fx := float64(x-stick) / float64(w-stick)
		half := 0.5 * (1 - fx)
		fy := float64(y) / float64(h)
		edge := 1.5 / float64(h) // a red border of about 1.5 pixels
		switch {
		case math.Abs(fy-0.5) > half:
			return noColor
		case math.Abs(fy-0.5) > half-edge, x-stick < 2:
			return juniorRed // the border, so the white stripes show on white
		case (x-stick)/max(w/7, 1)%2 == 0:
			return juniorRed
		default:
			return white
		}
	})
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

// garlandColors are the colors of the little flags of the garland.
var garlandColors = []color.NRGBA{juniorRed, sun, green, sky, orange, pink}

// newGarland returns a string of little triangular flags like the ones that
// hang over the streets during the Carnaval de Barranquilla. The string has
// color thread, so it shows over the sea and over a white card.
func newGarland(width, height float32, thread color.Color) fyne.CanvasObject {
	return newRaster(width, height, func(x, y, _, h int) color.Color {
		if y < max(h/9, 1) {
			return thread
		}
		i, t := x/h, float64(x%h)/float64(h) // flag number and position inside it
		depth := 1 - math.Abs(t-0.5)*2       // 1 in the middle, 0 at the sides
		if t > 0.08 && t < 0.92 && float64(y)/float64(h) <= depth {
			return garlandColors[i%len(garlandColors)]
		}
		return noColor
	})
}

// newSeaBanner returns the background of the welcome banner: the Caribbean
// sea, a sun on the right and two waves at the bottom, with rounded corners.
func newSeaBanner() fyne.CanvasObject {
	return canvas.NewRasterWithPixels(seaPixel)
}

// seaPixel returns the color of pixel (x, y) of a banner of w x h pixels.
func seaPixel(x, y, w, h int) color.Color {
	if !insideRounded(x, y, w, h, h/9) {
		return noColor
	}
	fx := float64(x) / float64(w)
	fy := float64(y) / float64(h)
	c := mix(caribe, sky, 0.8*fx+0.2*fy)

	// The sun: a disc with a soft glow around it. Distances are measured in
	// banner heights, so the sun stays round when the banner gets wider.
	d := math.Hypot((fx-0.9)*float64(w)/float64(h), fy-0.3)
	const disc, glow = 0.13, 0.42
	switch {
	case d < disc:
		c = sun
	case d < glow:
		t := (glow - d) / (glow - disc) // 1 next to the disc, 0 far away
		c = mix(c, sun, t*t*0.8)
	}

	// Two waves: below each sine line the sea gets a little lighter.
	for i, level := range []float64{0.76, 0.88} {
		wave := level + 0.03*math.Sin(fx*math.Pi*float64(5+2*i)+float64(i))
		if fy > wave {
			c = mix(c, white, 0.16)
		}
	}
	return c
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

// mix returns the color between a (t = 0) and b (t = 1).
func mix(a, b color.NRGBA, t float64) color.NRGBA {
	t = min(max(t, 0), 1)
	channel := func(p, q uint8) uint8 {
		return uint8(float64(p) + (float64(q)-float64(p))*t)
	}
	return color.NRGBA{R: channel(a.R, b.R), G: channel(a.G, b.G), B: channel(a.B, b.B), A: 0xFF}
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

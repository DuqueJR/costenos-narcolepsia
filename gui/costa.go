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
// Go, with no image files.

// noColor is a fully transparent pixel. It must be a color.NRGBA, like the
// other pixels: Fyne picks the kind of image from the type of the first
// pixel, and color.Transparent (a color.Alpha16) would make an image that
// only stores transparency and loses every color.
var noColor = color.NRGBA{}

// newVueltiaoBand returns a strip of width x height with the pattern of the
// sombrero vueltiao, the woven cane hat of the Colombian Caribbean coast.
func newVueltiaoBand(width, height float32) fyne.CanvasObject {
	r := canvas.NewRasterWithPixels(func(x, y, _, h int) color.Color {
		return vueltiaoPixel(x, y, h)
	})
	r.SetMinSize(fyne.NewSize(width, height))
	return r
}

// newAppIcon returns the icon of the app: a rounded square woven with the
// pattern of the sombrero vueltiao.
func newAppIcon(size float32) fyne.CanvasObject {
	r := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if !insideRounded(x, y, w, h, w/4) {
			return noColor
		}
		return vueltiaoPixel(x, y, h/3)
	})
	r.SetMinSize(fyne.NewSize(size, size))
	return r
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

// newSeaBanner returns the background of the welcome banner: the Caribbean
// sea, a sun on the right and two waves at the bottom, with rounded corners.
func newSeaBanner() fyne.CanvasObject {
	return canvas.NewRasterWithPixels(seaPixel)
}

// seaPixel returns the color of pixel (x, y) of a banner of w x h pixels.
// fx and fy are the position from 0 to 1, so the drawing does not depend on
// the size of the window.
func seaPixel(x, y, w, h int) color.Color {
	if !insideRounded(x, y, w, h, h/9) {
		return noColor
	}
	fx := float64(x) / float64(w)
	fy := float64(y) / float64(h)
	c := mix(caribe, sky, 0.8*fx+0.2*fy)

	// The sun: a disc with a soft glow around it. Distances are measured in
	// banner heights, so the sun stays round when the banner gets wider.
	d := math.Hypot((fx-0.86)*float64(w)/float64(h), fy-0.34)
	const disc, glow = 0.15, 0.45
	switch {
	case d < disc:
		c = sun
	case d < glow:
		t := (glow - d) / (glow - disc) // 1 next to the disc, 0 far away
		c = mix(c, sun, t*t*0.8)
	}

	// Two waves: below each sine line the sea gets a little lighter.
	for i, level := range []float64{0.74, 0.86} {
		wave := level + 0.035*math.Sin(fx*math.Pi*float64(5+2*i)+float64(i))
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

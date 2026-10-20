package main

import (
	"embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// iconFiles holds the icons of the GUI, from Phosphor Icons (MIT license,
// see icons/LICENSE). The embed package of the standard library puts them
// inside the binary at build time, so the app needs no extra files and no
// new Go dependency.
//
//go:embed icons/*.svg
var iconFiles embed.FS

// pictureFiles holds the pictures of the GUI, from Wikimedia Commons: the
// sombrero vueltiao and the flag of Junior. images/CREDITS.md says who made
// them and under which license.
//
//go:embed images/*.png
var pictureFiles embed.FS

// icon returns the icon called name, such as "bed", as a themed resource:
// Fyne paints it with the colors of the theme, like its own icons.
func icon(name string) fyne.Resource {
	data, err := iconFiles.ReadFile("icons/" + name + ".svg")
	if err != nil {
		// The names are written in the code and the files are checked when
		// building, so this only happens with a typo. The question mark
		// shows it on screen instead of stopping the app.
		return theme.QuestionIcon()
	}
	return theme.NewThemedResource(fyne.NewStaticResource(name+".svg", data))
}

// picture returns the picture called name, such as "sombrero", as a
// resource. Like icon, a typo shows a question mark instead of stopping
// the app.
func picture(name string) fyne.Resource {
	data, err := pictureFiles.ReadFile("images/" + name + ".png")
	if err != nil {
		return theme.QuestionIcon()
	}
	return fyne.NewStaticResource(name+".png", data)
}

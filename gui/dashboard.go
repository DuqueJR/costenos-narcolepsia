package main

import (
	"fmt"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// recentEpisodes is how many episodes the dashboard shows.
const recentEpisodes = 5

// newDashboardPage builds the Summary page of app a: a welcome banner, the
// big counters, the bed occupancy, the staff, the costeño saying of the
// day, quick actions and the last episodes. It is called "Resumen" and not
// "Inicio" because a specific name says what is inside.
func newDashboardPage(a *App) page {
	date := canvas.NewText("", muted)
	date.TextSize = 13

	beds := newMeter(0, 0, accent)
	bedsCount := newHeading("", theme.SizeNameHeadingText)
	bedsText := newCaption("")
	occupancy := newTitledPanel("Ocupación de camas",
		container.NewVBox(bedsCount, padded(2, beds), padded(2, bedsText)))

	doctors := newHeading("0", theme.SizeNameHeadingText)
	orderlies := newHeading("0", theme.SizeNameHeadingText)
	staff := newTitledPanel("Personal de turno", container.NewGridWithColumns(2,
		container.NewVBox(doctors, padded(2, newCaption("Doctores"))),
		container.NewVBox(orderlies, padded(2, newCaption("Camilleros"))),
	))

	saying := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	saying.Wrapping = fyne.TextWrapWord
	quote := newTitledPanel("Frase costeña del día", saying)

	quick := newTitledPanel("Acciones rápidas", container.NewGridWithColumns(2,
		newActionButton("Admitir paciente", icon("user-plus"), widget.HighImportance, a.showAdmitPatientForm),
		newActionButton("Agregar cuarto", icon("bed"), widget.MediumImportance, a.showAddRoomForm),
		newActionButton("Contratar doctor", icon("stethoscope"), widget.MediumImportance, a.showHireDoctorForm),
		newActionButton("Contratar camillero", icon("identification-badge"), widget.MediumImportance, a.showHireOrderlyForm),
	))

	recent := container.NewVBox()
	last := newTitledPanel("Últimos episodios", recent)

	a.addRefresher(func() {
		now := time.Now()
		date.Text = spanishDate(now)
		date.Refresh()
		saying.SetText("“" + sayingOfTheDay(now) + "”")

		used, total := 0, 0
		for _, r := range a.h.Rooms() {
			used += len(r.Occupants())
			total += r.Capacity()
		}
		bedsCount.SetText(fmt.Sprintf("%d / %d", used, total))
		beds.set(float64(used), float64(total), accent)
		bedsText.Text = fmt.Sprintf("Esperando cama en el pasillo: %d", len(a.h.PatientsInHallway()))
		bedsText.Refresh()

		doctors.SetText(fmt.Sprint(len(a.h.Doctors())))
		orderlies.SetText(fmt.Sprint(len(a.h.Staff()) - len(a.h.Doctors())))

		history := a.h.History() // a copy, so reversing it does not change the hospital
		slices.Reverse(history)
		recent.RemoveAll()
		if len(history) == 0 {
			recent.Add(padded(4, newCaption("Todo tranquilo, como domingo en Puerto Colombia.")))
		}
		for _, rec := range history[:min(len(history), recentEpisodes)] {
			line := widget.NewLabel(rec.Summary())
			// Without truncation a long line would widen the whole page.
			line.Truncation = fyne.TextTruncateEllipsis
			recent.Add(line)
		}
		recent.Refresh()
	})

	content := container.NewVBox(
		newBanner(date),
		newStatsStrip(a),
		container.NewGridWithColumns(3, occupancy, staff, quote),
		container.NewGridWithColumns(2, quick, last),
	)
	return page{
		section:  "Hospital",
		title:    "Resumen",
		subtitle: "¿Qué más pues? Así va el hospital hoy",
		icon:     icon("squares-four"),
		content:  container.NewVScroll(content),
	}
}

// newBanner returns the welcome card: today's date, a greeting and, on the
// right, a sombrero vueltiao. A thin vueltiao stripe closes it, as a
// signature: the costeño detail is there, but it does not shout.
func newBanner(date *canvas.Text) fyne.CanvasObject {
	title := canvas.NewText("¡Ajá, bienvenido al hospital!", ink)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	line := canvas.NewText("Aquí nadie se queda dormío sin que lo atiendan.", muted)
	line.TextSize = 14

	text := container.NewVBox(layout.NewSpacer(), date, title, line, layout.NewSpacer())
	hat := container.NewCenter(newSombrero(132, 86))
	top := container.New(layout.NewCustomPaddedLayout(4, 12, 8, 8), container.NewBorder(nil, nil, nil, hat, text))
	return newPanel(container.NewBorder(nil, newVueltiaoBand(10, 8), nil, nil, top))
}
